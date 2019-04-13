<template>
  <div class="data-container">
    <h2 class="title">Chores</h2>
    <div>{{ errorMessage }}</div>
    <div class="chore" v-for="chore in chores" v-bind:key="chore.chore_id">{{ chore.chore_id }}</div>
    <button v-on:click="newChore">New chore</button>
  </div>
</template>
<script>
import axios from "axios";
import dateFormat from "dateformat";

const user = "mizutani";
const now = new Date();
const today = dateFormat(now, "yyyy-mm-dd");

const choreData = {
  chores: [],
  errorMessage: ""
};

function newChore() {}

export default {
  data() {
    return choreData;
  },
  methods: {
    newChore: newChore
  }
};
</script>
<style lang="css" scoped>
</style>