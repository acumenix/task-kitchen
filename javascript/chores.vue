<template>
  <div class="chore-container">
    <h2 class="title">Chores</h2>
    {{ color }}
  </div>
</template>
<script>
const data = {
  color: "blue"
};

export default {
  data() {
    return data;
  }
};
</script>
<style lang="css" scoped>
</style>