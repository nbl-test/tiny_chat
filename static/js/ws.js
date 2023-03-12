// format:
let default_message = {
    from: "uuid",
    to: "uuid",
    sent_time: "xxxx",
    content: "content",
    // calculated
    is_public: true, //false,
    sent_by_me: true, // false
}
let default_model = {
    messages: [], // [message]
    username: ""
}
let default_error = {
    error: "error_content"
}
let default_identity = {
    name: "uuid"
}

// model
let model = {
    messages: [],
    username: ""
}

let message_callback = {
    on_set_username: null,
    on_add_msg: null,
    on_connection_lost: null,
    on_reset_msg: null,
}

function set_username(name) {
    model.username = name
    if (message_callback.on_set_username) {
        message_callback.on_set_username(name)
    }
}

function reset_messages() {
    Object.assign(model, default_model)

    console.log(model)
    console.log(default_model)
    if (message_callback.on_reset_msg) {
        message_callback.on_reset_msg()
    }
}

function add_message(message) {
    message.is_public = message.to?.length == 0;
    message.sent_by_me = message.from == username;
    model.messages.push(message)
    if (message_callback.on_add_msg) {
        message_callback.on_add_msg(message)
    }
}

function get_messages(skip_public, offset, count) {
    if (!offset) offset = 0;
    if (!count) count = 10;
    if (count + offset > model.messages.length) {
        count = model.messages.length - offset
    }
    if (!skip_public) {
        let start = model.messages.length - (offset+count)
        let end = model.messages.length - offset
        return model.messages.slice(start, end).reverse()
    }
    let ret = [];
    for (let idx = model.messages.length-1;idx >= 0;idx--) {
        if (!model.messages[idx].to) {
            // public
            continue
        }
        ret.push(model.messages[idx])
        if (ret.length == count) {
            break
        }
    }
    return ret
}

function connection_closed(event) {
    if (message_callback.on_connection_lost) {
        message_callback.on_connection_lost(event)
    }
}

// returns close function
// events must provides: {onopen, onmessage, onerror, onclose}
function init_ws(url, events) {
    const ws = new WebSocket(url);
    if (events) {
        let e = ()=>{}
        ws.onopen = (event) => {
            console.log(event);
            (events['onopen']||e)(event)
        };
        ws.onmessage = (event) => {
            console.log(event.data);
            (events['onmessage']||e)(event)
        }
        ws.onerror = (event) => {
            console.log(event);
            (events['onerror']||e)(event)
        }
        ws.onclose = (event) => {
            console.log(event);
            (events['onclose']||e)(event)
        }
    }
    return () => {
        ws.close();
    }
}

function init_message_ws() {
    const url = new URL('/chat', location.href);
    url.protocol = 'wss';
    return init_ws(url, {
        'onopen': (event)=>{
            reset_messages();
        },
        'onmessage': (event)=>{
            let msg = JSON.parse(event.data)
            if ('error' in msg) {
                alert('error occured: ' + msg)
                return
            }
            if ('name' in msg) {
                set_username(msg.name)
                return
            }
            if ('content' in msg && 'from' in msg) {
                add_message(msg)
                return
            }
        },
        'onerror': (event)=>{
            
        },
        'onclose': (event)=>{
            connection_closed(event)
        },
    })
}
