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
  this.name = json.type;
  this.max = options.max || BroTop.max;
  this.count = 0;
  this.items = [];
  this.allowScroll = true;

  this.template = BroTop.templates.collection;
  this.sidebarTemplate = BroTop.templates.collectionSidebar;

  this.id = "#" + this.name;

  if ($(this.id).length <= 0) {
    $("#wrapper .content").append(this.template(json));
    $("#sidebar ul.sidebar-menu").append(this.sidebarTemplate(json))
  }

}

Collection.prototype.Minimize = function() {
  var self = this;

  self.max = 5;

  BroTop.current = null;
  self.allowScroll = false;

  self.Cleanup();

  $(self.id).removeClass("full-screen")

  $(self.id).find(".display").off("scroll");
}

Collection.prototype.Maximise = function() {
  var self = this;

  BroTop.current = self.name;

  self.max = 1000;

  self.Cleanup();

  console.log(self.id)

  $(self.id).addClass("full-screen")

  $(self.id).find(".display").on("scroll", function(e) {
    var elem = $(e.currentTarget);

    console.log(elem)

    if (elem[0].scrollHeight - elem.scrollTop() == elem.outerHeight()) {
      self.allowScroll = true;
    } else {
      self.allowScroll = false;
    }
  });
}

Collection.prototype.Show = function() {
  $(this.id).show();
}

Collection.prototype.Hide = function() {
  $(this.id).hide();
}


Collection.prototype.Cleanup = function() {
  var self = this;

  while (self.count > self.max) {

    var evt = self.items.shift()

    if (evt) {
      evt.Remove();
    }

    this.count--;
  }
}

Collection.prototype.Scroll = function() {
  var self = this;
  var $item = $(self.id).find(".display");

  if (self.allowScroll) {
    $item.scrollTop($item[0].scrollHeight)
  }
}

Collection.prototype.Add = function(json) {
  var self = this;

  while (self.count >= self.max) {

    var evt = self.items.shift()

    if (evt) {
      evt.Remove();
    }

    this.count--;
  }

  var event = new Event(json);
  this.items.push(event);
  event.Render();
  this.count++;

}

var BroTop;

BroTop = {
  max: 5,
  collection: {},
  templates: {},
  Count: {},

  Graph: {
    Run: null,

    update: function(type) {
      if (BroTop.Count.hasOwnProperty(type)) {
        BroTop.Count[type]++;
      } else {
        BroTop.Count[type] = 1;
      }
    },
    send: function() {
      var time = (new Date).getTime();
      var data = [];
      var count = 0;

      for (item in BroTop.Count) {
        count += BroTop.Count[item];
        BroTop.Count[item] = 0;
      }

      BroTop.Graph.TS.append(time, count);
    },

    stop: function() {
      clearInterval(BroTop.Graph.Run);
    }
  },

  Init: function() {

    $(document).on("click", "ul.sidebar-menu li.item a", function(e) {
      e.preventDefault();
      var self = $(this).parent("li");
      var type = self.attr("data-type");

      $("ul.sidebar-menu li.item").removeClass("selected");
      self.addClass("selected");

      if (type === "all") {

        for (item in BroTop.collection) {
          var i = BroTop.collection[item]
          i.Minimize();
          i.Show();
        }
      
      } else {
        for (item in BroTop.collection) {
          var i = BroTop.collection[item]
          i.Minimize();
          i.Hide();
        }

        var c = BroTop.collection[type];

        console.log("GOT CLICK", self, type, c)
        c.Show();
        c.Maximise();
      }

    });

    BroTop.Graph.Run = setInterval(function() {
      BroTop.Graph.send();
    }, 1000);

    BroTop.Graph.TS = new TimeSeries();

    var chart = new SmoothieChart({
      grid: {
        fillStyle:'#14171b'
      }
    });

    chart.addTimeSeries(BroTop.Graph.TS, { 
      strokeStyle: '#e47078', 
      fillStyle: 'rgba(228,112,120,0.22)', 
      lineWidth: 4 
    });
    chart.streamTo(document.getElementById("main-chart"), 500);

    gotalk.handleNotification('event', function (event) {
      var json = JSON.parse(event)
      // console.log(json)

      if (json.hasOwnProperty("type")) {

        if (json.type === BroTop.current) {
          BroTop.collection[json.type].Scroll();
        }

        BroTop.Graph.update(json.type);

        if (BroTop.collection.hasOwnProperty(json.type)) {
          BroTop.collection[json.type].Add(json);
        } else {

          var collection = new Collection(json, {});

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
    this.max = max;
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

  var source = $("#sidebar-item").html();
  BroTop.templates.collectionSidebar = Handlebars.compile(source);

  BroTop.Init()
});
