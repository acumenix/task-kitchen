<template>
  <div class="data-container">
    <h2 class="title">Chores</h2>
    <div>{{ errorMessage }}</div>
    <div>
      <input
        class="new-chore edit"
        autofocus
        autocomplete="off"
        placeholder="Chore?"
        v-model="newChoreTitle"
        @keyup.enter="newChore"
        @keypress="enableSubmit"
      >
    </div>
    <div
      class="chore"
      v-for="chore in chores"
      v-bind:key="chore.chore_id"
      :class="{ editing: chore === editedChore, viewing: chore !== editedChore }"
    >
      <div class="view clearfix">
        <div @dblclick="editChore(chore)" class="title">
          <span v-if="chore.title !== ''">{{chore.title}}</span>
          <span v-else class="no-title">no title</span>
        </div>
      </div>
      <input
        class="edit"
        type="text"
        v-model="chore.title"
        v-chore-focus="chore == editedChore"
        @blur="doneEdit(chore)"
        @keyup.enter="doneEdit(chore)"
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

function newChore(event) {
  if (!isSubmitEnabled()) {
    return;
  }
  if (appData.newChoreTitle === "") {
    return;
  }

  axios
    .post(`/api/v1/${user}/${today}/chore`, { title: appData.newChoreTitle })
    .then(function(response) {
      appData.chores.push(response.data.results);
      appData.newChoreTitle = "";
      console.log(response);
    })
    .catch(function(error) {
      console.log(error);
      appData.errorMessage = "Error: " + error;
    });
}

function fetchChores() {
  axios
    .get(`/api/v1/${user}/${today}/chore`)
    .then(function(response) {
      // handle success
      if (response.data.results !== null) {
        appData.chores = appData.chores.concat(response.data.results);
      }
      console.log(response);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
      appData.errorMessage = "Error: " + error;
    });
}
fetchChores();

function updateChore(chore) {
  if (chore.prevTitle === undefined || chore.prevTitle === chore.title) {
    return;
  }
  console.log(chore);

  axios
    .put(`/api/v1/${user}/${today}/chore/${chore.chore_id}`, chore)
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

function editChore(chore) {
  this.editedChore = chore;
  chore.prevTitle = chore.title;
}

function doneEdit(chore) {
  if (!isSubmitEnabled()) {
    return;
  }

  if (!this.editedChore) {
    return;
  }
  this.editedChore = null;
  chore.title = chore.title.trim();
  updateChore(chore);
}

function enableSubmit(event) {
  appData.canSubmit = true;
  setTimeout(function() {
    appData.canSubmit = false;
  }, 200);
}
function isSubmitEnabled() {
  if (appData.canSubmit) {
    appData.canSubmit = false;
    return true;
  } else {
    return false;
  }
}

const appData = {
  chores: [],
  errorMessage: "",
  newChoreTitle: "",
  editedChore: null,
  canSubmit: false
};

export default {
  data() {
    return appData;
  },
  methods: {
    newChore: newChore,
    editChore: editChore,
    doneEdit: doneEdit,
    enableSubmit: enableSubmit
  },
  directives: {
    "chore-focus": function(el, binding) {
      if (binding.value) {
        el.focus();
      }
    }
  }
};
</script>
<style lang="css" scoped>
.chore {
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

.new-chore,
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