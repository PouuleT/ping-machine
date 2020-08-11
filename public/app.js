var app = new Vue({
  el: '#app',

  data: {
    message: 'Ping machine',
    pings: [],
    address: '8.8.8.8',
    count: 2,
    isLoading: false
  },

  methods: {
    ping: function() {
      this.isLoading = true;
      this.pings = [];
      axios.get("/ping/"+ this.count + "/times/"+ this.address)
        .then(this.onSuccess)
        .catch(error => {
          this.isLoading = false;
        });
    },

    onSuccess(response) {
      this.isLoading = false;
      this.pings = response.data;
    },
    getColorClass(ping) {
      return ping.response == 'OK' ? 'has-text-success' : 'has-text-danger';
    },
    getIcon(ping) {
      return ping.response == 'OK' ? 'fa-long-arrow-alt-right' : 'fa-times-circle';
    }
  }
});
