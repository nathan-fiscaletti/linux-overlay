const socket = new WebSocket(wsUrl);

socket.onopen = function(event) {
    console.log("Connected to WebSocket server");
};

socket.onmessage = function(event) {
    const data = JSON.parse(event.data);

    if (data.type === "key") {
        const key = document.getElementById(`${data.code}`);
        if (data.state === true) {
            if (!key.classList.contains("pressed")) {
                key.classList.add("pressed");
            }
        } else {
            key.classList.remove("pressed");
        }
    } else if (data.type === "mouse_movement") {
        const movement = document.getElementById("mouse_movement");
        movement.innerText = `${data.directions.join(", ")}`;
    }
};