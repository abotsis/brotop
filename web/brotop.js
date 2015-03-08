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


  if (!BroTop.Columns.hasOwnProperty(this.data.type)) {
    BroTop.Columns[this.data.type] = {}
  }

  var type = this.data.type;

  for (item in this.data.data) {
    var value = this.data.data[item].field

    if (!BroTop.Columns[type].hasOwnProperty(value)) {
      BroTop.Columns[type][value] = true
    }

    this.data.data[item].show = BroTop.Columns[type][value];
  }

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

  BroTop.Columns = {}
  $("table", self.id).find("th").show();
  $("table", self.id).find("td").show();

  $(".column-options", self.id).hide();
  $(self.id).find(".display").off("scroll");
}

Collection.prototype.Maximise = function() {
  var self = this;

  BroTop.current = self.name;

  self.max = 1000;

  self.Cleanup();

  $(self.id).addClass("full-screen")

  $(".display", self.id).on("scroll", function(e) {
    var elem = $(e.currentTarget);
    var t = $("table", e.currentTarget);

    if (elem[0].scrollHeight - elem.scrollTop() == elem.outerHeight()) {
      self.allowScroll = true;
    } else {
      self.allowScroll = false;
    }

  });

  $(".column-options", self.id).show();
  $(".display", self.id).trigger("scroll");
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

  var $item = $(".display", self.id);

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

  SetVersion: function() {
    if (BroTop.VersionRequest) {
      BroTop.VersionRequest.abort();
    }

    BroTop.VersionRequest = $.ajax({
      url: "/version",
      dataType: "json",
      type: "get",
      success: function(data) {
        $("span.version").html("v " + data.version)
      }
    });
  },

  max: 5,
  collection: {},
  templates: {},
  Count: {},

  Columns: {},

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
        fillStyle: '#14171b'
      }
    });

    chart.addTimeSeries(BroTop.Graph.TS, {
      strokeStyle: '#e47078',
      fillStyle: 'rgba(228,112,120,0.22)',
      lineWidth: 4
    });
    chart.streamTo(document.getElementById("main-chart"), 500);

    gotalk.handleNotification('event', function(event) {
      var json = JSON.parse(event)

      if (json.hasOwnProperty("type")) {

        BroTop.Graph.update(json.type);

        if (BroTop.collection.hasOwnProperty(json.type)) {
          BroTop.collection[json.type].Add(json);
        } else {

          var collection = new Collection(json, {});

          BroTop.collection[json.type] = collection;
        }

        BroTop.collection[json.type].Add(json);

        if (json.type === BroTop.current) {
          BroTop.collection[json.type].Scroll();
        }

      }


    });

    var s = gotalk.connection().on('open', function() {
      // ..
    });
  },

  ChangeMax: function(max) {
    this.max = max;
    for (collection in BroTop.collection) {
      var item = BroTop.collection[collection];
      item.max = max;
      item.Cleanup();
    }
  },

  ShowColumns: function(parent, type) {
    BroTop.Columns[parent][type] = true;
    var p = $("#" + parent);
    var t = $("table", p);

    t.find("th[data-name='"+type+"']").show()
    t.find("td[data-field='"+type+"']").show()
  },

  HideColumns: function(parent, type) {
    BroTop.Columns[parent][type] = false;
    var p = $("#" + parent);
    var t = $("table", p);

    t.find("th[data-name='"+type+"']").hide()
    t.find("td[data-field='"+type+"']").hide()
  }

}

jQuery(document).ready(function($) {

  BroTop.SetVersion();

  var source = $("#collection").html();
  BroTop.templates.collection = Handlebars.compile(source);

  var source = $("#event").html();
  BroTop.templates.event = Handlebars.compile(source);

  var source = $("#sidebar-item").html();
  BroTop.templates.collectionSidebar = Handlebars.compile(source);

  $(document).on("change", "input.column-checkbox", function() {
    var checked = $(this).is(":checked");
    var parent = $(this).attr("data-parent");
    var type = $(this).attr("data-type");

    if (!BroTop.Columns.hasOwnProperty(parent)) {
      BroTop.Columns[parent] = {} 
    }

    if (checked) {
      BroTop.ShowColumns(parent, type)
    } else {
      BroTop.HideColumns(parent, type)
    }
  });


  BroTop.Init()
});
