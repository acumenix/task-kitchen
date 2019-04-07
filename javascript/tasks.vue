<template>
  <div class="task-container">
    <h2 class="title">Tasks</h2>
    <div class="task" v-for="task in tasks" v-bind:key="task.task_id">{{ task.task_id }}</div>
    <button v-on:click="newTask">New task</button>
  </div>
</template>
<script>
import axios from "axios";

const user = "mizutani";
const taskData = {
  tasks: []
};

function newTask() {
  const today = "2019-04-07";

  axios
    .post(`/api/v1/$${user}/${today}/task`)
    .then(function(response) {
      // handle success
      taskData.tasks.push(response.data.results);
      console.log(response);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
    })
    .then(function() {
      // always executed
    });
}

function fetchTask() {
  const today = "2019-04-07";

  axios
    .get(`/api/v1/$${user}/${today}/task`)
    .then(function(response) {
      // handle success
      if (response.data.results !== undefined) {
        taskData.tasks = taskData.tasks.concat(response.data.results);
      }
      console.log(response);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
    })
    .then(function() {
      // always executed
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