var actionState = {
    name: "n/a",
    settings: {},
    globalSettings: {},
}



var debugInfo = {
    view: function (vnode) {
        return m("details", { class: "message" }, [
            m("summary", "Debug info"),
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

var actionUI = {
    "com.github.gebv.my-stream-deck-plugins.toggle-on-off":  onOffAction,
    "com.github.gebv.my-stream-deck-plugins.dosomething1":  testAction
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
            document.body,
            m("p", 'no registered Properoty Inspector for ' +actionState.name+ ' action'),
        )
        return
    }
    m.mount(
        document.body,
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
