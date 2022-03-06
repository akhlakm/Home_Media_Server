<script>
    import { onMount } from "svelte";
    import App from "./app";
    import { Canvas } from "./canvas";
    import Sidebar from "./Sidebar.svelte";

    onMount(() => {
        Canvas.LoadImage($App.url, $App.cvStatic);
        Canvas.SetSize($App.cvDynamic, Canvas.GetSize($App.cvStatic));
    });

    function handleCanvasClick(e) {
        e.stopPropagation();
        e.preventDefault();
        if (e.ctrlKey) {
            //if ctrl key is pressed, add new text
            var newtext = Canvas.GetNewText($App.texts, e.clientX, e.clientY);
            App.AddText(newtext);
            Canvas.SetSize($App.cvDynamic, Canvas.GetSize($App.cvStatic));
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
            document.getElementById("value").focus();
        } else {
            // if ctrl key is not pressed
        }
    }

    function handleKeyUp(e) {
        let keyCode = e.keyCode;
        let chrCode = keyCode - 48 * Math.floor(keyCode / 48);
        let chr = String.fromCharCode(96 <= keyCode ? chrCode : keyCode);
        console.log(chr, e);
    }
</script>

<div class="container">
    <div id="canvas-area">
        <canvas
            class="canvas"
            id={$App.cvStatic}
            style="z-Index: 0;"
            on:click={handleCanvasClick}
            on:keyup={handleKeyUp}
        />
        <canvas
            class="canvas"
            id={$App.cvDynamic}
            style="z-Index: 1;"
            on:click={handleCanvasClick}
            on:keyup={handleKeyUp}
        />
    </div>
    <div id="sidebar-area">
        <Sidebar />
    </div>
</div>

<style>
    .container {
        background-color: #414141;
        position: relative;
        display: flex;
        flex-flow: row nowrap;
        justify-content: center;
        align-items: flex-start;
        height: 100vh;
    }

    .container::-webkit-scrollbar {
        display: none;
    }

    #canvas-area {
        flex: 4 4 0;
        overflow-y: scroll;
        max-height: 100vh;
        position: relative;
        height: 100vh;
    }

    #sidebar-area {
        flex: 1 1 0;
        overflow-y: auto;
        max-height: 100vh;
        margin: auto;
    }

    canvas {
        position: absolute;
        left: 0;
        top: 0;
    }
</style>
