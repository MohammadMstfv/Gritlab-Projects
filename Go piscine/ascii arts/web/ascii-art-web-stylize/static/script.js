const form = document.getElementById("ascii-form");
const text = document.getElementById("text");
const banner = document.getElementById("banner");
const result = document.getElementById("result");
const liveResult = document.getElementById("live-result");
const generateButton = document.getElementById("generate-button");
const renderStatus = document.getElementById("render-status");
const resultStatus = document.getElementById("result-status");
const characterCount = document.getElementById("character-count");
const copyButton = document.getElementById("copy-button");
const saveButton = document.getElementById("save-button");
const exampleButton = document.getElementById("example-button");
const clearButton = document.getElementById("clear-button");
let renderTimer;
let previewRequest;
let generateRequest;
let hasGeneratedOutput = false;

function updateCount() {
    const count = text.value.length;
    characterCount.textContent = count + (count === 1 ? " character" : " characters");
    clearButton.disabled = count === 0;
}

function updateActionState() {
    copyButton.disabled = !hasGeneratedOutput;
    saveButton.disabled = !hasGeneratedOutput;
}

async function requestOutput(signal) {
    const response = await fetch(form.action, {
        method: "POST",
        body: new URLSearchParams({ text: text.value, banner: banner.value }),
        signal: signal
    });
    const output = await response.text();
    if (!response.ok) throw new Error(output);
    return output;
}

async function renderPreview() {
    updateCount();
    if (!text.value) {
        if (previewRequest) previewRequest.abort();
        liveResult.textContent = "Preview updates as you type.";
        liveResult.classList.add("is-empty");
        renderStatus.textContent = "Waiting for text";
        return;
    }

    if (previewRequest) previewRequest.abort();
    previewRequest = new AbortController();
    renderStatus.textContent = "Updating preview...";
    try {
        const output = await requestOutput(previewRequest.signal);
        liveResult.textContent = output;
        liveResult.classList.remove("is-empty");
        renderStatus.textContent = "Up to date";
    } catch (error) {
        if (error.name === "AbortError") return;
        liveResult.textContent = error.message || "Could not reach the server. Please try again.";
        liveResult.classList.add("is-empty");
        renderStatus.textContent = "Preview unavailable";
    }
}

async function generateOutput() {
    updateCount();
    if (!text.value) {
        text.reportValidity();
        return;
    }

    if (generateRequest) generateRequest.abort();
    generateRequest = new AbortController();
    resultStatus.textContent = "Generating...";
    try {
        const output = await requestOutput(generateRequest.signal);
        result.textContent = output;
        result.classList.remove("is-empty");
        hasGeneratedOutput = true;
        resultStatus.textContent = "Generated successfully";
        updateActionState();
    } catch (error) {
        if (error.name === "AbortError") return;
        result.textContent = error.message || "Could not reach the server. Please try again.";
        result.classList.add("is-empty");
        hasGeneratedOutput = false;
        resultStatus.textContent = "Could not generate art";
        updateActionState();
    }
}

function scheduleRender() {
    clearTimeout(renderTimer);
    if (hasGeneratedOutput) {
        resultStatus.textContent = "Input changed — generate again to refresh";
    }
    renderTimer = setTimeout(renderPreview, 140);
}

generateButton.addEventListener("click", function () {
    clearTimeout(renderTimer);
    generateOutput();
});
text.addEventListener("input", scheduleRender);
banner.addEventListener("change", renderPreview);
exampleButton.addEventListener("click", function () {
    text.value = "Hello, world!";
    text.focus();
    updateCount();
    if (hasGeneratedOutput) {
        resultStatus.textContent = "Input changed — generate again to refresh";
    }
    clearTimeout(renderTimer);
    renderPreview();
});
clearButton.addEventListener("click", function () {
    text.value = "";
    text.focus();
    updateCount();
    if (previewRequest) previewRequest.abort();
    liveResult.textContent = "Preview updates as you type.";
    liveResult.classList.add("is-empty");
    renderStatus.textContent = "Waiting for text";
    if (hasGeneratedOutput) {
        resultStatus.textContent = "Input changed — generate again to refresh";
    }
});
copyButton.addEventListener("click", async function () {
    try {
        await navigator.clipboard.writeText(result.textContent);
        resultStatus.textContent = "Copied to clipboard";
    } catch (error) {
        resultStatus.textContent = "Clipboard unavailable — select the art to copy";
    }
});
saveButton.addEventListener("click", function () {
    const lines = result.textContent.replace(/\n$/, "").split("\n");
    const padding = 32;
    const lineHeight = 20;
    const font = "16px monospace";
    const canvas = document.createElement("canvas");
    const context = canvas.getContext("2d");
    if (!context) {
        resultStatus.textContent = "Could not create the image";
        return;
    }
    context.font = font;
    const width = Math.max(...lines.map(function (line) { return context.measureText(line).width; }));
    canvas.width = Math.ceil(width + padding * 2);
    canvas.height = lines.length * lineHeight + padding * 2;
    context.fillStyle = "#182126";
    context.fillRect(0, 0, canvas.width, canvas.height);
    context.font = font;
    context.fillStyle = "#e8b84d";
    lines.forEach(function (line, index) {
        context.fillText(line, padding, padding + (index + 1) * lineHeight);
    });
    canvas.toBlob(function (blob) {
        if (!blob) {
            resultStatus.textContent = "Could not create the image";
            return;
        }
        const link = document.createElement("a");
        link.download = "ascii-art.png";
        link.href = URL.createObjectURL(blob);
        link.click();
        setTimeout(function () { URL.revokeObjectURL(link.href); }, 1000);
        resultStatus.textContent = "Image saved";
    }, "image/png");
});
updateActionState();
