<template>
  <div class="app-container">
    <h2 class="title">Fix Today's Task</h2>
  </div>
</template>
<script>
import axios from "axios";
import dateFormat from "dateformat";

const user = "mizutani";
const now = new Date();
const today = dateFormat(now, "yyyy-mm-dd");

const appData = {
  report: null
};

function getReport() {
  axios
    .get(`/api/v1/${user}/${today}`)
    .then(function(response) {
      // handle success
      if (response.data.results !== null) {
        appData.report = response.data.results;
      }
      console.log(response.data.results);
    })
    .catch(function(error) {
      // handle error
      console.log(error);
      appData.errorMessage = "Error: " + error;
    });
}

getReport();

export default {
  data() {
    return {};
  }
};
</script>
<style lang="css" scoped>
div.app-container {
  width: 640px;
  margin: 25px auto;
  padding: 15px;
  background-color: #fff;
  border: 1px solid #bbb;
}

h2.title {
  text-align: center;
  font-size: 24px;
}
</style>