<script>
    import { onMount } from "svelte";
    import App from "./app";
    import { Canvas } from "./canvas";
    import Sidebar from "./Sidebar.svelte";

    onMount(() => {
        Canvas.LoadImage($App.url, $App.cvStatic, (img) => {
            Canvas.SetSize($App.cvDynamic, Canvas.GetSize($App.cvStatic));
        });
    });

    function handleCanvasClick(e) {
        e.stopPropagation();
        e.preventDefault();
        var canvasarea = document.getElementById("canvas-area");
        if (e.ctrlKey) {
            //if ctrl key is pressed, add new text
            var newtext = Canvas.GetNewText(
                $App.texts,
                canvasarea.scrollLeft + e.clientX,
                canvasarea.scrollTop + e.clientY
            );
            App.AddText(newtext);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
            var inp = document.getElementById("value");
            inp.focus();
            setTimeout(() => {
                inp.setSelectionRange(0, inp.value.length);
                inp.select();
            }, 100);
        } else {
            // if ctrl key is not pressed, select the text
            var coord = {
                x: canvasarea.scrollLeft + e.clientX,
                y: canvasarea.scrollTop + e.clientY,
            };

            // a text hit?
            var hit = Canvas.GetTextHit($App.texts, coord);
            if (hit > -1) {
                App.SelectText(hit);
            }
        }
    }

    let isDragging = false;
    let startCoord = null;

    // drag start
    function handleMouseDown(e) {
        var canvasarea = document.getElementById("canvas-area");
        startCoord = {
            x: canvasarea.scrollLeft + e.clientX,
            y: canvasarea.scrollTop + e.clientY,
        };
        var hit = Canvas.GetTextHit($App.texts, startCoord);
        if (hit > -1) {
            isDragging = true;
            App.SelectText(hit);
        }
    }

    // drag end
    function handleMouseUp(e) {
        isDragging = false;
    }

    // normal move or drag move
    function handleMouseMove(e) {
        var canvasarea = document.getElementById("canvas-area");
        var coord = {
            x: canvasarea.scrollLeft + e.clientX,
            y: canvasarea.scrollTop + e.clientY,
        };
        if (isDragging) {
            var dx = coord.x - startCoord.x;
            var dy = coord.y - startCoord.y;
            App.UpdatePosition(dx, dy);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
            startCoord = coord;
            App.SetInfo(dx + " " + dy);
        } else {
            // a text hit, show info or cursor position
            var hit = Canvas.GetTextHit($App.texts, coord);
            if (hit > -1) {
                App.SetInfo("Text " + hit);
            } else {
                App.SetInfo(coord.x + " " + coord.y);
            }
        }
    }
</script>

<div class="container">
    <div id="canvas-area">
        <canvas class="canvas" id={$App.cvStatic} style="z-Index: 0;" />
        <canvas
            class="canvas"
            id={$App.cvDynamic}
            style="z-Index: 1;"
            on:click={handleCanvasClick}
            on:mousemove={handleMouseMove}
            on:mousedown={handleMouseDown}
            on:mouseup={handleMouseUp}
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
