const socket = new WebSocket(wsUrl);

socket.onopen = function(event) {
    console.log("Connected to WebSocket server");
};

socket.onmessage = function (event) {
    const data = JSON.parse(event.data);
    const mouseElement = document.getElementById("mouse_movement");
    let currentMouseDirection = null;

    if (data.type === "key") {
        handleKeyInput(data);
    } else if (data.type === "mouse_movement") {
        handleMouseMovement(data);
    }

    function handleKeyInput(data) {
        const { code, state } = data;

        // Mouse buttons
        if ([272, 273, 274].includes(code)) {
            mouseElement.className = state ? getMouseButtonClass(code) : "mouse-default";
            return
        }

        // Keyboard keys
        const key = document.getElementById(`${code}`);
        if (state) {
            key.classList.add("pressed");
        } else {
            key.classList.remove("pressed");
        }
    }

    function getMouseButtonClass(code) {
        switch (code) {
            case 272: return "mouse-lmb";
            case 273: return "mouse-rmb";
            case 274: return "mouse-mmb";
            default: return "mouse-default";
        }
    }

    function handleMouseMovement(data) {
        const newDirection = data.directions[0];
        if (newDirection !== currentMouseDirection) {
            mouseElement.innerHTML = `<img src="mousemove_${data.directions.join(", ")}.png" style="width: 50%;">`;
            currentMouseDirection = newDirection;
        }
    }
};