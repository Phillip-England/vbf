let store = {}
let events = []

function signal(initialValue) {
    let value = initialValue;
    let subscribers = [];

    function get() {
        return value;
    }

    function set(newValue) {
        if (value !== newValue) {
        value = newValue;
        subscribers.forEach(callback => callback(value));
        }
    }

    function subscribe(callback) {
        subscribers.push(callback);
        return () => {
        subscribers = subscribers.filter(sub => sub !== callback);
        };
    }

    return { get, set, subscribe };

}

function qs(selector, node=null) {
    if (node) {
        return node.querySelector(selector)
    } else {
        return document.querySelector(selector)
    }
}

function qsa(selector, node=null) {
    if (node) {
        return node.querySelectorAll(selector)
    } else {
        return document.querySelectorAll(selector)
    }
}