package login_service_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gorm.io/gorm"

	common "test.com/project-common"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/jwts"
	"test.com/project-common/tms"
	"test.com/project-grpc/account"
	"test.com/project-grpc/user/login"
	"test.com/project-user/config"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/data/organization"
	"test.com/project-user/internal/database/gorms"
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/rpc"
	"test.com/project-user/pkg/model"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache                repo.Cache
	memberRepo           repo.MemberRepo
	organizationRepo     repo.OrganizationRepo
	accountServiceClient account.AccountServiceClient
}

func New() *LoginService {
	return &LoginService{
		cache:                dao.Rc,
		memberRepo:           dao.NewMemberDao(),
		organizationRepo:     dao.NewOrganizationDao(),
		accountServiceClient: rpc.AccountServiceClient,
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	// 1.获取参数
	mobile := msg.Mobile
	// 2.校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	// 3.生成验证码（随机4位1000-9999或者6位100000-999999）
	code := "123456"
	// 4.调用短信平台（三方 放入go协程中执行 接口可以快速响应）
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功，发送短信")
		// redis 假设后续缓存可能存在mysql当中，也可能存在mongo当中 也可能存在memcache当中
		// 5.存储验证码 redis当中 过期时间15分钟
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := ls.cache.Put(c, model.RegisterRedisKey+mobile, code, 15*time.Minute)
		if err != nil {
			zap.L().Info(fmt.Sprintf("验证码存入redis出错,cause by: %v \n", err))
		}
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	c := context.Background()
	// 1.可以校验参数
	// 2.校验验证码
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+msg.Mobile)
	if err == redis.Nil {
		return nil, errs.GrpcError(model.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("Register redis get error", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	// 3.校验业务逻辑（邮箱是否被注册 账号是否被注册 手机号是否被注册）
	exist, err := ls.memberRepo.GetMemberByEmail(c, msg.Email)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	exist, err = ls.memberRepo.GetMemberByAccount(c, msg.Name)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	exist, err = ls.memberRepo.GetMemberByMobile(c, msg.Mobile)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}
	// 4.执行业务 将数据存入member表 生成一个数据 存入组织表 organization
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	err = gorms.GetDB().Transaction(func(tx *gorm.DB) error {
		err = ls.memberRepo.SaveMember(tx, c, mem)
		if err != nil {
			zap.L().Error("Register db SaveMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		// 存入组织
		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			// TODO: Personal 个人组织标识写死为 1，应该改为配置化
			Personal: model.Personal,
			// TODO: Avatar 头像 URL 是写死的，应该改为配置化或使用默认头像
			Avatar: "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = ls.organizationRepo.SaveOrganization(tx, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		// 创建用户账户并绑定默认角色
		// 修改说明：新增账户创建逻辑，用户注册时自动创建账户记录
		// 背景：项目需要实现基于角色的权限管理，用户注册时需要自动创建账户并绑定默认角色
		// 实现：调用project-project服务的gRPC接口创建账户，设置IsOwner=1表示创建者为所有者
		orgCode := fmt.Sprintf("%d", org.Id)
		accountCode := fmt.Sprintf("%d", mem.Id)
		_, err = ls.accountServiceClient.SaveAccount(c, &account.AccountReqMessage{
			MemberId:         mem.Id,
			OrganizationCode: orgCode,
			AccountCode:      accountCode,
			Authorize:        "1",
			IsOwner:          1,
			Name:             mem.Name,
			Mobile:           mem.Mobile,
			Email:            mem.Email,
			Avatar:           mem.Avatar,
		})
		if err != nil {
			zap.L().Error("Register SaveAccount error", zap.Error(err))
		}
		return nil
	})
	// 5. 返回
	return &login.RegisterResponse{}, err
}

// Login 用户登录
// 修改说明：增加手机验证码登录方式支持
// 背景：原系统仅支持账号密码登录，现增加手机号+验证码登录方式以提升用户体验
// 修改点：
// 1. 判断请求中是否同时包含Mobile和Captcha，如果是则走验证码登录流程
// 2. 验证码登录流程：从Redis获取验证码->校验验证码->通过手机号查询用户
// 3. 账号密码登录流程保持不变
func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	c := context.Background()

	var mem *member.Member
	var err error

	if msg.Mobile != "" && msg.Captcha != "" {
		redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+msg.Mobile)
		if err == redis.Nil {
			return nil, errs.GrpcError(model.NoLegalCaptcha)
		}
		if err != nil {
			zap.L().Error("Login get captcha error", zap.Error(err))
			return nil, errs.GrpcError(model.NoLegalCaptcha)
		}
		if redisCode != msg.Captcha {
			return nil, errs.GrpcError(model.CaptchaError)
		}
		mem, err = ls.memberRepo.FindMemberByMobile(c, msg.Mobile)
		if err != nil {
			zap.L().Error("Login db FindMemberByMobile error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if mem == nil {
			return nil, errs.GrpcError(model.AccountAndPwdError)
		}
	} else {
		pwd := encrypts.Md5(msg.Password)
		mem, err = ls.memberRepo.FindMember(c, msg.Account, pwd)
		if err != nil {
			zap.L().Error("Login db FindMember error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if mem == nil {
			return nil, errs.GrpcError(model.AccountAndPwdError)
		}
	}
	memMsg := &login.MemberMessage{}
	err = copier.Copy(memMsg, mem)
	memMsg.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey)
	memMsg.LastLoginTime = tms.FormatByMill(mem.LastLoginTime)
	memMsg.CreateTime = tms.FormatByMill(mem.CreateTime)
	// 2.根据用户id查组织
	orgs, err := ls.organizationRepo.FindOrganizationByMemId(c, mem.Id)
	if err != nil {
		zap.L().Error("Login db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, v := range orgsMessage {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
		v.OwnerCode = memMsg.Code
		o := organization.ToMap(orgs)[v.Id]
		v.CreateTime = tms.FormatByMill(o.CreateTime)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	// 3.用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.C.JwtConfig.AccessExp*3600*24) * time.Second
	rExp := time.Duration(config.C.JwtConfig.RefreshExp*3600*24) * time.Second
	token := jwts.CreateToken(memIdStr, exp, config.C.JwtConfig.AccessSecret, rExp, config.C.JwtConfig.RefreshSecret, msg.Ip)
	// 可以给token做加密处理 增加安全性
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		AccessTokenExp: token.AccessExp,
		TokenType:      "bearer",
	}
	// 放入缓存 member orgs
	go func() {
		marshal, _ := json.Marshal(mem)
		ls.cache.Put(c, model.Member+"::"+memIdStr, string(marshal), exp)
		orgsJson, _ := json.Marshal(orgs)
		ls.cache.Put(c, model.MemberOrganization+"::"+memIdStr, string(orgsJson), exp)
	}()
	return &login.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}

func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	parseToken, err := jwts.ParseToken(token, config.C.JwtConfig.AccessSecret, msg.Ip)
	if err != nil {
		zap.L().Error("Login  TokenVerify error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	// 从缓存中查询 如果没有 直接返回认证失败
	memJson, err := ls.cache.Get(context.Background(), model.Member+"::"+parseToken)
	if err != nil {
		zap.L().Error("TokenVerify cache get member error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	if memJson == "" {
		zap.L().Error("TokenVerify cache get member expire")
		return nil, errs.GrpcError(model.NoLogin)
	}
	memberById := &member.Member{}
	json.Unmarshal([]byte(memJson), memberById)
	// 数据库查询 优化点 登录之后 应该把用户信息缓存起来
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)

	orgsJson, err := ls.cache.Get(context.Background(), model.MemberOrganization+"::"+parseToken)
	if err != nil {
		zap.L().Error("TokenVerify cache get organization error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	if orgsJson == "" {
		zap.L().Error("TokenVerify cache get organization expire")
		return nil, errs.GrpcError(model.NoLogin)
	}
	var orgs []*organization.Organization
	json.Unmarshal([]byte(orgsJson), &orgs)

	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	memMsg.CreateTime = tms.FormatByMill(memberById.CreateTime)
	return &login.LoginResponse{Member: memMsg}, nil
}

func (ls *LoginService) Logout(ctx context.Context, msg *login.LogoutMessage) (*login.LogoutResponse, error) {
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	parseToken, err := jwts.ParseToken(token, config.C.JwtConfig.AccessSecret, "")
	if err != nil {
		zap.L().Error("Logout ParseToken error", zap.Error(err))
		return &login.LogoutResponse{Success: false}, nil
	}
	c := context.Background()
	err = ls.cache.Del(c, model.Member+"::"+parseToken)
	if err != nil {
		zap.L().Error("Logout Del member cache error", zap.Error(err))
	}
	err = ls.cache.Del(c, model.MemberOrganization+"::"+parseToken)
	if err != nil {
		zap.L().Error("Logout Del organization cache error", zap.Error(err))
	}
	zap.L().Info("User logout success, token: " + parseToken)
	return &login.LogoutResponse{Success: true}, nil
}

func (l *LoginService) MyOrgList(ctx context.Context, msg *login.UserMessage) (*login.OrgListResponse, error) {
	memId := msg.MemId
	orgs, err := l.organizationRepo.FindOrganizationByMemId(ctx, memId)
	if err != nil {
		zap.L().Error("MyOrgList FindOrganizationByMemId err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	orgMap := organization.ToMap(orgs)
	for _, org := range orgsMessage {
		org.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey)
		org.OwnerCode, _ = encrypts.EncryptInt64(orgMap[org.Id].MemberId, model.AESKey)
		org.CreateTime = tms.FormatByMill(orgMap[org.Id].CreateTime)
	}
	return &login.OrgListResponse{OrganizationList: orgsMessage}, nil
}

func (ls *LoginService) SaveOrganization(ctx context.Context, msg *login.OrganizationMessage) (*login.OrganizationMessage, error) {
	org := &organization.Organization{
		Name:       msg.Name,
		MemberId:   msg.MemberId,
		CreateTime: time.Now().UnixMilli(),
		Personal:   0,
		Address:    msg.Address,
	}
	err := gorms.GetDB().Transaction(func(tx *gorm.DB) error {
		return ls.organizationRepo.SaveOrganization(tx, ctx, org)
	})
	if err != nil {
		zap.L().Error("SaveOrganization db err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	msg.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey)
	msg.Id = org.Id
	msg.CreateTime = tms.FormatByMill(org.CreateTime)
	msg.OwnerCode = msg.Code
	return msg, nil
}

func (ls *LoginService) UpdateOrganization(ctx context.Context, msg *login.OrganizationMessage) (*login.OrganizationMessage, error) {
	id := encrypts.DecryptNoErr(msg.Code)
	if id == 0 {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	org := &organization.Organization{
		Id:      id,
		Name:    msg.Name,
		Address: msg.Address,
	}
	if err := ls.organizationRepo.UpdateOrganization(ctx, org); err != nil {
		zap.L().Error("UpdateOrganization db err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return msg, nil
}

func (ls *LoginService) FindMemInfoById(ctx context.Context, msg *login.UserMessage) (*login.MemberMessage, error) {
	memberById, err := ls.memberRepo.FindMemberById(context.Background(), msg.MemId)
	if err != nil {
		zap.L().Error("TokenVerify db FindMemberById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)
	orgs, err := ls.organizationRepo.FindOrganizationByMemId(context.Background(), memberById.Id)
	if err != nil {
		zap.L().Error("TokenVerify db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	memMsg.CreateTime = tms.FormatByMill(memberById.CreateTime)
	return memMsg, nil
}

func (ls *LoginService) FindMemInfoByIds(ctx context.Context, msg *login.UserMessage) (*login.MemberMessageList, error) {
	zap.L().Info("FindMemInfoByIds debug", zap.Any("MIds", msg.MIds))
	memberList, err := ls.memberRepo.FindMemberByIds(context.Background(), msg.MIds)
	zap.L().Info("FindMemInfoByIds debug memberList", zap.Any("memberList", memberList), zap.Int("len", len(memberList)))
	if err != nil {
		zap.L().Error("FindMemInfoByIds db memberRepo.FindMemberByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if memberList == nil || len(memberList) <= 0 {
		return &login.MemberMessageList{List: nil}, nil
	}
	mMap := make(map[int64]*member.Member)
	for _, v := range memberList {
		mMap[v.Id] = v
	}
	var memMsgs []*login.MemberMessage
	copier.Copy(&memMsgs, memberList)
	for _, v := range memMsgs {
		m := mMap[v.Id]
		v.CreateTime = tms.FormatByMill(m.CreateTime)
		v.Code = encrypts.EncryptNoErr(v.Id)
	}

	return &login.MemberMessageList{List: memMsgs}, nil
}

func (ls *LoginService) EditPersonal(ctx context.Context, msg *login.MemberMessage) (*login.MemberMessage, error) {
	if msg.Id == 0 {
		return nil, errs.GrpcError(model.NoLogin)
	}
	fields := map[string]interface{}{}
	if msg.Avatar != "" {
		fields["avatar"] = msg.Avatar
	}
	// uploadAvatar 只传 Avatar，不传 Name/Description，所以用 msg.Name != "" 判断是否为 editPersonal 调用
	if msg.Name != "" {
		fields["name"] = msg.Name
	}
	// description 由 editPersonal 调用时传入，uploadAvatar 不传，用 Name 非空作为 editPersonal 的标志
	// 当 Name 非空时（editPersonal 场景），同步更新 description（允许清空）
	if msg.Name != "" {
		fields["description"] = msg.Description
	}
	if len(fields) == 0 {
		return msg, nil
	}
	if err := ls.memberRepo.UpdateMemberFields(context.Background(), msg.Id, fields); err != nil {
		zap.L().Error("EditPersonal UpdateMemberFields error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return msg, nil
}
