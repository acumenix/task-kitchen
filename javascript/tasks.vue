<template>
  <div class="data-container">
    <h2 class="title">Tasks</h2>
    <div>{{ errorMessage }}</div>
    <div class="task" v-for="task in tasks" v-bind:key="task.task_id">{{ task.task_id }}</div>
    <button v-on:click="newTask">New task</button>
  </div>
</template>
<script>
import axios from "axios";
import dateFormat from "dateformat";

const user = "mizutani";
const now = new Date();
const today = dateFormat(now, "yyyy-mm-dd");

const taskData = {
  tasks: [],
  errorMessage: ""
};

function newTask() {
  axios
    .post(`/api/v1/${user}/${today}/task`)
    .then(function(response) {
      taskData.tasks.push(response.data.results);
      console.log(response);
    })
    .catch(function(error) {
      console.log(error);
      taskData.errorMessage = "Error: " + error;
    });
}

function fetchTask() {
  axios
    .get(`/api/v1/${user}/${today}/task`)
    .then(function(response) {
      // handle success
      if (response.data.results !== null) {
        taskData.tasks = taskData.tasks.concat(response.data.results);
      }
      console.log(response);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
      taskData.errorMessage = "Error: " + error;
    });
}

fetchTask();

export default {
  data() {
    return taskData;
  },
  methods: {
    newTask: newTask
  }
};
</script>
<style lang="css" scoped>
</style>