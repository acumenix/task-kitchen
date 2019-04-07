<template>
  <div class="g-content">
    <tasks></tasks>
    <chores></chores>
  </div>
</template>
<script>
import tasks from "./tasks.vue";
import chores from "./chores.vue";

export default {
  components: {
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
.g-content {
  background-color: #eee;
}
.title {
  font-size: 24px;
  font-weight: bold;
}
</style>