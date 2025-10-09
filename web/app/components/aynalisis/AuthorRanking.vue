<template>
  <div class="mb-4">
    <q-table dense :title="i18n.t('contributors')" color="primary" :rows="rows" :columns="columns"
      v-model:pagination="pagination" hide-pagination row-key="nickname">
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props"
            class="text-purple text-weight-bold text-lg bg-yellow-2">
            {{ col.label }}
          </q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td key="nickname" :props="props">
            <div class="text-blue text-weight-bold text-lg">
              {{ props.row.nickname }}
            </div>
          </q-td>
          <q-td key="name" :props="props">
            <div v-html="props.row.name?.replace(/,/g, '<br>')"></div>
          </q-td>
          <q-td key="email" :props="props">
            <div v-html="props.row.email?.replace(/,/g, '<br>')"></div>
          </q-td>
          <q-td key="commits" :props="props">
            <div class="text-yellow text-weight-bold text-lg">
              {{ props.row.commits }}
            </div>
          </q-td>
          <q-td key="additions" :props="props">
            <div class="text-green text-weight-bold text-lg">
              {{ props.row.additions }}
            </div>
          </q-td>
          <q-td key="deletions" :props="props">
            <div class="text-red text-weight-bold text-lg">
              {{ props.row.deletions }}
            </div>
          </q-td>
          <q-td key="effectives" :props="props">
            <div class="text-teal text-weight-bold text-lg">
              {{ props.row.effectives }}
            </div>
          </q-td>
          <q-td key="projects" :props="props">
            <div class="text-purple text-weight-bold text-lg">
              {{ props.row.projects }}
            </div>
          </q-td>
        </q-tr>
      </template>

      <template v-slot:bottom>
        <div class="flex flex-row nowrap" >
          <div class="mr-1">{{ $t("timeRange") }}:</div>
          <div v-if="since" color="primary" text-color="white">
            {{ since?.split?.(" ")?.[0] }}
          </div>
          <div class="mx-1">~</div>
          <div v-if="until" color="primary" text-color="white">
            {{ until?.split?.(" ")?.[0] }}
          </div>
        </div>
      </template>
    </q-table>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import type { QTableColumn } from "quasar";

const i18n = useI18n();

const props = defineProps({
  authors: {
    type: Array,
    required: true
  },
  since: {
    type: String,
    required: true
  },
  until: {
    type: String,
    required: true
  },
  sortBy: {
    type: String,
  },
  sortDirection: {
    type: String,
  }
})

const pagination = ref({
  rowsPerPage: 0,
  sortBy: props.sortBy,
  descending: props.sortDirection === "desc",
});

const columns: QTableColumn<any>[] = [
  {
    name: "nickname",
    label: i18n.t("name"),
    align: "center",
    field: "nickname",
    sortable: true,
  },
  { name: "name", label: i18n.t("nickname"), align: "center", field: "name" },
  { name: "email", label: i18n.t("email"), align: "center", field: "email" },
  {
    name: "commits",
    label: i18n.t("commits"),
    align: "center",
    field: "commits",
    sortable: true
  },
  {
    name: "additions",
    label: i18n.t("additions"),
    align: "center",
    field: "additions",
    sortable: true,
  },
  {
    name: "deletions",
    label: i18n.t("deletions"),
    align: "center",
    field: "deletions",
    sortable: true,
  },
  {
    name: "effectives",
    label: i18n.t("effectives"),
    align: "center",
    field: "effectives",
    sortable: true,
  },
  {
    name: "projects",
    label: i18n.t("projects"),
    align: "center",
    field: "projects",
    sortable: true,
  },
] as const;

const rows = computed(() => props.authors);

</script>
