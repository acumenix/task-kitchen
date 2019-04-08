<template>
  <div id="main">
    <controller></controller>
    <tasks></tasks>
    <chores></chores>
  </div>
</template>
<script>
import controller from "./controller.vue";
import tasks from "./tasks.vue";
import chores from "./chores.vue";

export default {
  components: {
    controller,
    tasks,
    chores
  },
  data() {
    return {
      greeting: "Hello"
    };
  }
};
</script>
<style lang="css" scoped>
</style>