var guid = (function() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
               .toString(16)
               .substring(1);
  }
  return function() {
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
           s4() + '-' + s4() + s4() + s4();
  };
})();

var Event;

function Event(json) {
  this.data = json; 
  this.parent = $("#" + json.type);
  this.id = guid();
  this.data.id = this.id;
  this.template = BroTop.templates.event;
}

Event.prototype.Render = function() {
  this.parent.find("tbody").append(this.template(this.data))
}

Event.prototype.Remove = function() {
  $("#" + this.id).remove();
}

var Collection;

function Collection(json, options) {
  this.name = json.name;
  this.max = options.max || 1;
  this.count = 0;
  this.items = [];

  this.template = BroTop.templates.collection;

  this.id = "#" + this.name;

  if ($(this.id).length <= 0) {
    $("body").append(this.template(json));
  }
}

Collection.prototype.Show = function() {
  $(this.id).show();
}

Collection.prototype.Hide = function() {
  $(this.id).hide();
}


Collection.prototype.Cleanup = function() {
  var self = this;

  while (self.count >= self.max) {
    var evt = self.items.pop()
    if (evt) {
      evt.Remove();
      this.count--;
    }
  }
}

Collection.prototype.Add = function(json) {
  var self = this;

  self.Cleanup();

  var event = new Event(json);
  this.items.push(event);
  event.Render();
  this.count++;

}

var BroTop;

BroTop = {

  collection: {},
  templates: {},

  Init: function() {
    gotalk.handleNotification('event', function (event) {
      var json = JSON.parse(event)
      // console.log(json)

      if (json.hasOwnProperty("type")) {

        if (BroTop.collection.hasOwnProperty(json.type)) {
          BroTop.collection[json.type].Add(json);
        } else {

          var collection = new Collection(json, {
            max: 1
          });

          BroTop.collection[json.type] = collection;
        }

        BroTop.collection[json.type].Add(json);
      }


    });

    gotalk.connect('ws://'+document.location.host+'/gotalk', function (err, s) {
      if (err) return console.error(err);
      // s is a gotalk.Sock
    });
  },

  ChangeMax: function(max) {
    for (collection in BroTop.collection) {
      var item = BroTop.collection[collection];
      item.max = max;
      item.Cleanup();
    }
  }

}

jQuery(document).ready(function($) {
  var source = $("#collection").html();
  BroTop.templates.collection = Handlebars.compile(source);

  var source = $("#event").html();
  BroTop.templates.event = Handlebars.compile(source);

  BroTop.Init()
});
