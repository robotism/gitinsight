<template>
    <div class="w-full flex-col content-center">
        <h6 class="w-full content-center">{{ $t("conditionFilter") }}</h6>

        <div class="w-full flex flex-col">
            <q-list bordered class="rounded-borders">
                <q-expansion-item class="w-full group1" expand-separator default-opened icon="📅"
                    :label="$t('timeRange')">
                    <q-card class="w-full px-2 pb-2 max-h-[20vh] overflow-y-auto">
                        <q-input dense v-model="range.since" type="date" clearable :prefix="$t('since')" />
                        <q-input dense v-model="range.until" type="date" clearable :prefix="$t('until')" />
                        <q-input dense v-model="range.geEffective" flat clearable type="number"
                            :prefix="$t('effectives') + '>='" />
                        <q-input dense v-model="range.leEffective" flat clearable type="number"
                            :prefix="$t('effectives') + '<='" />
                    </q-card>
                </q-expansion-item>

                <q-expansion-item class="w-full group3 nowrap" dense expand-separator default-opened icon="👥"
                    :label="$t('contributors')" header-class="text-purple">
                    <q-card class="w-full pl-2 pb-2 max-h-[20vh] overflow-y-auto">
                        <q-option-group class="w-full text-nowrap flex-wrap" dense v-model="authorSelections"
                            :options="authorOptions" color="purple" type="checkbox">
                            <template v-slot:label="opt">
                                <span class="text-purple text-[8px]">{{ opt.label }}</span>
                            </template>
                        </q-option-group>
                    </q-card>
                </q-expansion-item>

                <q-expansion-item class="w-full group1" dense expand-separator default-opened icon="🌿"
                    :label="$t('repos')">
                    <q-card class="w-full pl-2 pb-2 max-h-[20vh] overflow-y-auto">
                        <q-option-group class="w-full" dense v-model="repoSelections" :options="repoOptions"
                            color="green" type="checkbox">
                            <template v-slot:label="opt">
                                <span class="text-teal text-[8px]">{{ opt.label }}</span>
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
import DateInput from "~/components/tool/DatePicker.vue";

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
    since: "",
    until: "",
    leEffective: "",
    geEffective: "",
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
                " (" +
                commitCount +
                " commits, " +
                authorCount +
                " authors)",
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

watch(
    [range, repoSelections, authorSelections],
    () => {
        emits("update:filter", {
            since: range.value?.since || "",
            until: range.value?.until || "",
            leEffective: range.value?.leEffective || "",
            geEffective: range.value?.geEffective || "",
            repos: repoSelections.value.filter((item: any) => !!item),
            authors: authorSelections.value.filter((item: any) => !!item),
        });
    },
    { deep: true }
);

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
