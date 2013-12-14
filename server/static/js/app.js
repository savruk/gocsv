App = Ember.Application.create();

App.Router.map(function() {
  // this.resource("index")
  this.resource("files")
  this.resource("about")
});


App.FilesRoute = Ember.Route.extend({
    model: function(){
        return $.getJSON("/files", function (items) {
            return items;
        });
    }
});
