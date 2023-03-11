// format:
let default_message = {
    from: "uuid",
    to: "uuid",
    sent_time: "xxxx",
    content: "content",
}
let default_model = {
    messages: [], // [message]
}
let default_error = {
    error: "error_content"
}

// model
let model = {
    messages: [],
}

function reset_messages() {
    Object.assign(model, default_model)

    console.log(model)
    console.log(default_model)
}

function add_message(message) {
    model.messages.push(message)
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

// export {add_message, get_messages, reset_messages}