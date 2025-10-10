<template>
    <div class="w-full flex-col content-center">
        <h6 class="w-full content-center">{{ $t("conditionFilter") }}</h6>

        <div class="w-full flex flex-col">
            <q-list bordered class="rounded-borders">
                <q-expansion-item class="w-full group1" expand-separator default-opened icon="📅"
                    :label="$t('timeRange')" header-class="bg-primary">
                    <q-card class="w-full px-2 pb-2 max-h-[20vh] overflow-y-auto">
                        <div class="flex flex-row flex-wrap">
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('')" :label="$t('all')" />
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('today')" :label="$t('today')" />
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('yesterday')"
                                :label="$t('yesterday')" />

                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('thisWeek')"
                                :label="$t('weekThis')" />
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('lastWeek')"
                                :label="$t('weekLast')" />
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('beforeLastWeek')"
                                :label="$t('weekBeforeLast')" />

                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('thisMonth')"
                                :label="$t('monthThis')" />
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('lastMonth')"
                                :label="$t('monthLast')" />
                            <q-btn class="mx-1" flat size="xs" @click="setDateRange('beforeLastMonth')"
                                :label="$t('monthBeforeLast')" />
                        </div>
                        <q-input dense v-model="range.since" type="date" clearable :prefix="$t('since')" min="2025-01-01" />
                        <q-input dense v-model="range.until" type="date" clearable :prefix="$t('until')" min="2025-01-01" />
                        <q-input dense v-model="range.geEffective" flat clearable type="number"
                            :prefix="$t('effectives') + '>='" />
                        <q-input dense v-model="range.leEffective" flat clearable type="number"
                            :prefix="$t('effectives') + '<='" />
                    </q-card>
                </q-expansion-item>

                <q-expansion-item class="w-full group3 nowrap" dense expand-separator default-opened icon="👥"
                    :label="$t('contributors')" header-class="bg-primary">
                    <q-card class="w-full pl-2 pb-2 max-h-[20vh] overflow-y-auto">
                        <q-option-group class="w-full text-nowrap flex-wrap" dense v-model="authorSelections"
                            :options="authorOptions" type="checkbox">
                            <template v-slot:label="opt">
                                <span class="text-[10px]" :style="{ color: hashColor(opt.label) }">{{ opt.label
                                    }}</span>
                            </template>
                        </q-option-group>
                    </q-card>
                </q-expansion-item>

                <q-expansion-item class="w-full group1" dense expand-separator default-opened icon="🌿"
                    :label="$t('repos')" header-class="bg-primary">
                    <q-card class="w-full pl-2 pb-2 max-h-[20vh] overflow-y-auto">
                        <q-option-group class="w-full" dense v-model="repoSelections" :options="repoOptions"
                            type="checkbox">
                            <template v-slot:label="opt">
                                <span class="text-[8px]"
                                    :style="{ color: hashColor('#' + opt.label?.split('-')?.[0]?.toLowerCase?.() + '-') }">
                                    {{ opt.label }}
                                </span>
                            </template>
                        </q-option-group>
                    </q-card>
                </q-expansion-item>

            </q-list>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { watch } from "vue";

import moment from "moment";

const { hashColor } = useColor();

const api = useApi();
const i18n = useI18n();

const props = defineProps({
    filter: {
        type: Object,
        required: true,
    },
});
const emits = defineEmits(["update:filter"]);

const range: Ref<any> = ref({
    since: props.filter?.since || "",
    until: props.filter?.until || "",
    leEffective: props.filter?.leEffective || "",
    geEffective: props.filter?.geEffective || "",
});

const repos: Ref<any[]> = ref([]);
const authors: Ref<any[]> = ref([]);


const repoSelections: Ref<string[]> = ref([]);
const repoOptions: ComputedRef<any[]> = computed(() => {
    const repoUrls = repos.value.map((item: any) => {
        const repoName = item.repoUrl.split("/").pop().replace(".git", "");
        const commitCount = item.commits?.split?.(",")?.length || item.commits || 0;
        const authorCount = item.nicknames?.split?.(",")?.length || 0;
        return {
            label:
                repoName +
                " (" + authorCount + " " + i18n.t("contributors") + ")",
            value: item.repoUrl,
        };
    });
    const list = [];
    for (const item of repoUrls || []) {
        if (list.findIndex((i: any) => i.value === item.value) < 0) {
            list.push(item);
        }
    }
    list.sort((a: any, b: any) => a.label.localeCompare(b.label));
    const fix = [
        {
            label: i18n.t("all"),
            value: "",
        },
    ];
    fix.push(...list);
    return fix;
});

const isRepoUpdating = ref(false);
watch(
    () => repoSelections.value,
    (newVal: string[], oldVal: string[]) => {
        if (isRepoUpdating.value) return;
        isRepoUpdating.value = true;

        const isNewAll = newVal.includes("");
        const isOldAll = oldVal.includes("");
        if (isOldAll && !isNewAll) {
            // 点击了取消全选，清空
            repoSelections.value = [];
        } else if (!isOldAll && isNewAll) {
            // 点击了全选，选中所有
            repoSelections.value = repoOptions.value.map((item: any) => item.value);
        } else if (newVal.length == repoOptions.value.length - 1 && !isNewAll) {
            // 只有全选没打勾，自动打勾
            repoSelections.value.splice(0, 0, "");
        } else if (newVal.length < repoOptions.value.length && isNewAll) {
            // 全选打勾了但是部分取消了，自动取消全选打勾
            repoSelections.value = repoSelections.value.filter((item: any) => !!item);
        }
        // 延迟一帧再解除锁，确保内部修改不会触发新的 watcher
        nextTick(() => {
            isRepoUpdating.value = false;
        });
    },
    { deep: true }
);

const authorSelections: Ref<string[]> = ref([]);
const authorOptions: ComputedRef<any[]> = computed(() => {
    const list = authors.value.map((item: any) => {
        return {
            label: item.nickname,
            value: item.nickname,
        };
    });
    list.sort((a, b) => a.label.localeCompare(b.label));
    const fix = [
        {
            label: i18n.t("all"),
            value: "",
        },
    ];
    fix.push(...list);
    return fix;
});
const isAuthorUpdating = ref(false);
watch(
    () => authorSelections.value,
    (newVal: string[], oldVal: string[]) => {
        if (isAuthorUpdating.value) return;
        isAuthorUpdating.value = true;

        const isNewAll = newVal.includes("");
        const isOldAll = oldVal.includes("");
        if (isOldAll && !isNewAll) {
            // 点击了取消全选，清空
            authorSelections.value = [];
        } else if (!isOldAll && isNewAll) {
            // 点击了全选，选中所有
            authorSelections.value = authorOptions.value.map(
                (item: any) => item.value
            );
        } else if (newVal.length == authorOptions.value.length - 1 && !isNewAll) {
            // 只有全选没打勾，自动打勾
            authorSelections.value.splice(0, 0, "");
        } else if (newVal.length < authorOptions.value.length && isNewAll) {
            // 全选打勾了但是部分取消了，自动取消全选打勾
            authorSelections.value = authorSelections.value.filter(
                (item: any) => !!item
            );
        }
        // 延迟一帧再解除锁，确保内部修改不会触发新的 watcher
        nextTick(() => {
            isAuthorUpdating.value = false;
        });
    },
    { deep: true }
);

const debounce = ref()
watch(
    [range, repoSelections, authorSelections],
    () => {
        clearTimeout(debounce.value)
        debounce.value = setTimeout(() => {
            emits("update:filter", {
                since: range.value?.since || "",
                until: range.value?.until || "",
                leEffective: range.value?.leEffective || "",
                geEffective: range.value?.geEffective || "",
                repos: repoSelections.value.filter((item: any) => !!item),
                authors: authorSelections.value.filter((item: any) => !!item),
            });
        }, 300)
    },
    { deep: true }
);

function setDateRange(type: string) {
    const [since, until] = getDateRange(type)
    range.value.since = since || range.value.since
    range.value.until = until || range.value.until
}

function getDateRange(type: string) {
    const FORMAT = "YYYY-MM-DD"
    const m = moment()

    switch (type) {
        case "today":
            return [m.startOf("day").format(FORMAT), m.endOf("day").format(FORMAT)]
        case "yesterday":
            return [
                m.clone().subtract(1, "day").startOf("day").format(FORMAT),
                m.clone().subtract(1, "day").endOf("day").format(FORMAT)
            ]
        case "thisWeek":
            return [m.startOf("week").format(FORMAT), m.endOf("week").format(FORMAT)]
        case "lastWeek":
            return [
                m.clone().subtract(1, "week").startOf("week").format(FORMAT),
                m.clone().subtract(1, "week").endOf("week").format(FORMAT)
            ]
        case "beforeLastWeek":
            return [
                m.clone().subtract(2, "week").startOf("week").format(FORMAT),
                m.clone().subtract(2, "week").endOf("week").format(FORMAT)
            ]
        case "thisMonth":
            return [m.startOf("month").format(FORMAT), m.endOf("month").format(FORMAT)]
        case "lastMonth":
            return [
                m.clone().subtract(1, "month").startOf("month").format(FORMAT),
                m.clone().subtract(1, "month").endOf("month").format(FORMAT)
            ]
        case "beforeLastMonth":
            return [
                m.clone().subtract(2, "month").startOf("month").format(FORMAT),
                m.clone().subtract(2, "month").endOf("month").format(FORMAT)
            ]
        default:
            return []
    }
}

const getData = () => {
    api.getRepoBranches().then((res: any) => {
        const list: any[] = [];
        for (const item of res?.data || []) {
            const exists = list.find((i: any) => i.repoUrl === item.repoUrl);
            if (exists) {
                exists.commits = Math.max(exists.commits, item.commits);
                exists.additions = Math.max(exists.additions, item.additions);
                exists.deletions = Math.max(exists.deletions, item.deletions);
                exists.effectives = Math.max(exists.effectives, item.effectives);
                exists.nicknames = Array.from(
                    new Set(
                        exists.nicknames
                            .split(",")
                            .concat(item.nicknames.split(","))
                            .filter((item: any) => !!item)
                    )
                ).join(",");
                exists.branchName = Array.from(
                    new Set(
                        exists.branchName
                            .split(",")
                            .concat(item.branchName.split(","))
                            .filter((item: any) => !!item)
                    )
                ).join(",");
            } else {
                list.push({
                    repoUrl: item.repoUrl,
                    branchName: item.branchName,
                    commits: item.commits,
                    nicknames: item.nicknames,
                    additions: item.additions,
                    deletions: item.deletions,
                    effectives: item.effectives,
                });
            }
        }
        repos.value = list;
    });
    api.getContributors().then((res: any) => {
        authors.value = res.data;
    });
};

onMounted(() => {
    getData();
});
</script>

<style scoped>
:deep(.q-item__section--side > .q-icon) {
    font-size: 14px;
    margin-bottom: 4px;
}

:deep(.group3 .q-option-group) {
    display: flex;
}

:deep(.q-field__prefix) {
    color: var(--q-color-primary) !important;
}
</style>
