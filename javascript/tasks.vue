<template>
  <div class="data-container">
    <h2 class="title">Tasks</h2>
    <div>{{ errorMessage }}</div>
    <div>
      <input
        class="new-task edit"
        autofocus
        autocomplete="off"
        placeholder="What is your task?"
        v-model="newTaskTitle"
        @keyup.enter="newTask"
        @keypress="enableSubmit"
      >
    </div>
    <div
      class="task"
      v-for="task in tasks"
      v-bind:key="task.task_id"
      :class="{ editing: task === editedTask, viewing: task !== editedTask }"
    >
      <div class="view clearfix">
        <div @dblclick="editTask(task)" class="title">
          <span v-if="task.title !== ''">{{task.title}}</span>
          <span v-else class="no-title">no title</span>
        </div>
        <div
          class="tomatos"
          @click.left.prevent="incrementTomato(task)"
          @click.right.prevent="decrimentTomato(task)"
          @mouseenter="editTomato(task)"
          @mouseleave="saveTomato(task)"
        >
          <span class="tomato" v-for="n in task.tomato_num" v-bind:key="n">🍅</span>
        </div>
      </div>
      <input
        class="edit"
        type="text"
        v-model="task.title"
        v-task-focus="task == editedTask"
        @blur="doneEdit(task)"
        @keyup.enter="doneEdit(task)"
        @keypress="enableSubmit"
      >
    </div>
  </div>
</template>
<script>
import axios from "axios";
import dateFormat from "dateformat";

const user = "mizutani";
const now = new Date();
const today = dateFormat(now, "yyyy-mm-dd");

function newTask(event) {
  if (!isSubmitEnabled()) {
    return;
  }
  if (appData.newTaskTitle === "") {
    return;
  }

  axios
    .post(`/api/v1/${user}/${today}/task`, { title: appData.newTaskTitle })
    .then(function(response) {
      appData.tasks.push(response.data.results);
      appData.newTaskTitle = "";
      console.log(response);
    })
    .catch(function(error) {
      console.log(error);
      appData.errorMessage = "Error: " + error;
    });
}

function fetchTask() {
  axios
    .get(`/api/v1/${user}/${today}/task`)
    .then(function(response) {
      // handle success
      if (response.data.results !== null) {
        appData.tasks = appData.tasks.concat(response.data.results);
      }
      console.log(response);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
      appData.errorMessage = "Error: " + error;
    });
}
fetchTask();

function updateTask(task) {
  if (
    (task.prevTitle === undefined || task.prevTitle === task.title) &&
    (task.prevTomatoNum === undefined || task.prevTomatoNum === task.tomato_num)
  ) {
    return;
  }
  console.log(task);

  axios
    .put(`/api/v1/${user}/${today}/task/${task.task_id}`, task)
    .then(function(response) {
      // handle success
      console.log(response);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
      appData.errorMessage = "Error: " + error;
    });
}

function editTask(task) {
  this.editedTask = task;
  task.prevTitle = task.title;
}

function doneEdit(task) {
  if (!isSubmitEnabled()) {
    return;
  }

  if (!this.editedTask) {
    return;
  }
  this.editedTask = null;
  task.title = task.title.trim();
  updateTask(task);
}

// ==========================================
// Tomato Control
function editTomato(task) {
  task.prevTomatoNum = task.tomato_num;
}

function saveTomato(task) {
  updateTask(task);
  task.prevTomatoNum = task.tomato_num;
}

function cancelEdit(task) {
  this.editTask = null;
  task.title = task.prevTitle;
}

function incrementTomato(task) {
  task.tomato_num++;
}

function decrimentTomato(task) {
  if (task.tomato_num <= 0) {
    return;
  }
  task.tomato_num--;
}

function enableSubmit(event) {
  appData.canSubmit = true;
  console.log("enable", event);
  setTimeout(function() {
    appData.canSubmit = false;
  }, 200);
}
function isSubmitEnabled() {
  console.log("submit", appData.canSubmit);
  if (appData.canSubmit) {
    appData.canSubmit = false;
    return true;
  } else {
    return false;
  }
}

const appData = {
  tasks: [],
  errorMessage: "",
  newTaskTitle: "",
  editedTask: null,
  canSubmit: false
};

export default {
  data() {
    return appData;
  },
  methods: {
    newTask: newTask,
    editTask: editTask,
    doneEdit: doneEdit,
    cancelEdit: cancelEdit,
    editTomato: editTomato,
    saveTomato: saveTomato,
    incrementTomato: incrementTomato,
    decrimentTomato: decrimentTomato,
    enableSubmit: enableSubmit
  },
  directives: {
    "task-focus": function(el, binding) {
      if (binding.value) {
        el.focus();
      }
    }
  }
};
</script>
<style lang="css" scoped>
.task {
  margin: 10px;
  color: #888;
  position: relative;
  font-size: 24px;
  border-bottom: 1px solid #ededed;
}

div.editing .view {
  display: none;
}

div.viewing .edit {
  display: none;
}

.view {
  position: relative;
  margin: 0;
  font-size: 24px;
  line-height: 1.4em;
  border: 0;
  padding: 6px;
  box-sizing: border-box;
}

.new-task,
.edit {
  position: relative;
  margin: 0;
  width: 100%;
  font-size: 24px;
  line-height: 1.4em;
  color: inherit;

  border: 0;
  padding: 6px;
  border: 1px solid #999;
  box-shadow: inset 0 -1px 5px 0 rgba(0, 0, 0, 0.2);
  box-sizing: border-box;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

input::-webkit-input-placeholder,
input::-moz-placeholder,
input::input-placeholder {
  font-weight: 300;
  color: #e6e6e6;
}

span.no-title {
  color: #eee;
}

div.tomatos {
  float: right;
  position: relative;
  margin: 0;
  font-size: 24px;
  line-height: 1.4em;
}
span.tomato {
  margin: 2px;
}
div.title {
  float: left;
}

.clearfix::after {
  content: "";
  display: block;
  clear: both;
}
</style>