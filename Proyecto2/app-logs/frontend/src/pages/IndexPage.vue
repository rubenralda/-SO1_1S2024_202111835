<template>
  <q-page class="row justify-center">
    <div class="q-pa-lg col-md-8 col-xs-12 col-sm-12" style="max-width: 1000px">
      <q-list dense bordered padding class="rounded-borders">
        <q-item clickable v-ripple v-for="(item, index) in logs" :key="index">
          <q-item-section> {{ item.fecha }} - {{ item.name }} </q-item-section>
        </q-item>
      </q-list>
    </div>
    <div class="q-pa-lg col-md-4">
      <q-btn
        align="between"
        class="btn-fixed-width"
        color="secondary"
        label="Update Logs"
        icon="flight_takeoff"
        :onclick="getLogs"
      />
    </div>
  </q-page>
</template>

<script setup>
import { ref } from "vue";

const logs = ref([{ name: "dfdf", fecha: "2024" }]);

function getLogs() {
  fetch("https://backend-szlvbetafa-uc.a.run.app/api/logs")
    .then((respone) => respone.json())
    .then((data) => logs.value = data);
}
</script>
