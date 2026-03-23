<!-- TODO: 日程日历功能暂未实现，需要新建数据库表 project_events 等 -->
<template>
    <div class="common-calendar">
        <p class="muted text-center" style="padding: 40px;">日程日历功能暂未实现</p>
    </div>
</template>

<!--
<template>
    <calendar></calendar>
</template>

<script>
import _ from 'lodash'
import pagination from "@/mixins/pagination";
import calendar from "@/components/project/calendar";

export default {
    name: "CommonCalendar",
    components: {
        calendar
    },
    mixins: [pagination],
    data() {
        return {
            keyword: '',
            selectedKeys: ['0'],
        }
    },
    watch: {
        keyword() {
            this.search();
        },
    },
    created() {
    },
    methods: {
    }
}
</script>
-->

<script>
export default {
    name: "CommonCalendar",
}
</script>
