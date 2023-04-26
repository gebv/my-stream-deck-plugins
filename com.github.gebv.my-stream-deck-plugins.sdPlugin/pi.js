var actionState = {
    context: "n/a",
    name: "n/a",
    settings: {},
    globalSettings: {},
}



var debugInfo = {
    view: function (vnode) {
        return m("div", [
            m("p", "Context ID:" + actionState.context),
            m("p", "Plugin settings:" + JSON.stringify(actionState.settings)),
            m("p", "Plugin globalSettings:" + JSON.stringify(actionState.globalSettings)),
        ])
    }
}

var onOffAction = {
    view: function (vnode) {
        return m("div", [
            m("h3", "On Off action"),
            m("p", "action name: " + actionState.name),
            m("button", {onclick: function() {
                $PI.sendToPlugin({key: "some key", value: "some value"});
            }}, "sendToPlugin"),
            m("button", {onclick: function() {
                // $PI.openUrl("getting-started.html")
                // window.xtWindow = window.open('getting-started.html', "Getting started");
                window.open('getting-started.html', "Getting started");
            }}, "openUrl"),
        ])
    }
}

var testAction = {
    view: function (vnode) {
        return m("div", [
            m("h3", "test Actions"),
            m("p", "action name: " + actionState.name),
        ])
    }
}

var memInfoAction = {
    selectedSkin: "",
    oninit: function(vnode) {
        memInfoAction.selectedSkin = actionState.settings.selectedSkin || "cpu_usage_percent"
    },
    availableSkins: [
        {name: "CPU Usage Percent", id: "cpu_usage_percent"},
        {name: "CPU Usage Percent (High-Performance Cores)", id: "cpu_usage_percent_hpc"},
        {name: "CPU Usage Percent (High-Efficiency Cores)", id: "cpu_usage_percent_hec"},
        {name: "Memory Usage Percent", id: "mem_usage_percent"},
        {name: "Memory Total", id: "mem_total"},
        {name: "Memory Free", id: "mem_free"},
    ],
    view: function() {
        return m("div", [
            m("h3", "Mem Info Settings"),
            // <div class="sdpi-item">
            m("div", {class: "sdpi-item"}, [
                // <div class="sdpi-item-label">Select</div>
                m("div", {class: "sdpi-item-label"}, "Skin"),
                // <select class="sdpi-item-value select"
                m("select", {
                    class: "sdpi-item-value select",
                    value: memInfoAction.selectedSkin,
                    onchange: e => {
                        memInfoAction.selectedSkin = e.target.value
                        actionState.settings.selectedSkin = e.target.value
                        console.log("selected skin", e.target)

                        $PI.setSettings(actionState.settings)
                    }},
                    memInfoAction.availableSkins.map(e => m("option", {value: e.id}, e.name))
                )
            ]),
            m("div", {class: "sdpi-item"}, [
                m("div.sdpi-item-label", "Debug Info"),
                m(debugInfo),
            ])
        ])
    }
}

var actionUI = {
    "com.github.gebv.my-stream-deck-plugins.toggle-on-off":  onOffAction,
    "com.github.gebv.my-stream-deck-plugins.dosomething1":  testAction,
    "com.github.gebv.my-stream-deck-plugins.mem-info":  memInfoAction,
}

// var itemExample = {
//     view: function (vnode) {
//         // vnode.actionState.data
//         return m("div", { class: "sdpi-item" }, [
//             m("div", { class: "sdpi-item-label" }, "Value From Settings (key=abc) " + actionState.inputValue),
//             m("input[type=text]", {
//                 oninput: (e) => { actionState.inputValue = e.target.value },
//                 class: "sdpi-item-value",
//             }),
//         ])
//     }
// }

function handleMyCustomEvent(e) {
    console.log("received my custom event", e)
}

document.addEventListener('myCustomEvent', handleMyCustomEvent);

function initUI() {
    console.log("mithril", m)
    if (!actionUI.hasOwnProperty(actionState.name)) {
        m.mount(
            document.getElementById("app"),
            m("p", 'no registered Properoty Inspector for ' +actionState.name+ ' action'),
        )
        return
    }
    m.mount(
        document.getElementById("app"),
        actionUI[actionState.name],
    )
}

console.log("$PI", $PI)

$PI.onDidReceiveGlobalSettings(e => {
    console.log('Received global settings', e);
    actionState.globalSettings = e.payload.settings
    m.redraw()
})

$PI.onConnected(e => {
    // TODO
    // $PI.loadLocalization('../../../');

    console.log("onConnected", e)

    const { actionInfo, appInfo, connection, messageType, port, uuid } = e;
    const { action, payload, context } = actionInfo;
    actionState.context = context
    const { settings } = payload;
    actionState.name = action
    actionState.settings = settings
    // $PI.setSettings(value);

    initUI()

    $PI.onDidReceiveSettings(action, e => {
        console.log('Received action settings', action, e);
        actionState.settings = e.payload.settings
        m.redraw()
    })

    $PI.onSendToPropertyInspector(action, e => {
        console.log('onSendToPropertyInspector', e);
    });

    $PI.getSettings()
    $PI.getGlobalSettings()

    console.log('Property Inspector connected', e);
});
