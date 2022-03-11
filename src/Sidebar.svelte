<script>
    import App from "./app";
    import { Canvas } from "./canvas";

    function handleValueChange(e) {
        App.UpdateText(e.target.value);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }

    function handleSelectText(e) {
        App.SelectText(parseInt(e.target.value));
    }

    function handleSelectColor(e, value) {
        if (e.ctrlKey) {
        } else {
            App.SelectColor(value);
        }
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }

    function handleSelectFont(e) {
        App.UpdateFont(e.target.value);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }

    function handleUpdateSize(e) {
        App.UpdateSize(e.target.value);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }
    function handleUpdateFitwidth(e) {
        App.UpdateFitwidth(e.target.value);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }
    function handleUpdateTextAlign(e) {
        App.UpdateAlign(e.target.value);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }

    function handleSaveClick(e) {
        if ($App.hash == null) {
            e.target.value = "No Hash!";
            return;
        } else {
            e.target.value = "Saving...";
            Canvas.Redraw($App.cvStatic, $App.texts, false);
            Canvas.GetBlob($App.cvStatic, (blob) => {
                var formData = new FormData();
                formData.append("file", blob);

                var request = new XMLHttpRequest();
                request.onload = () => {
                    e.target.value = "Saved";
                };

                request.open("POST", "/api/caption/" + $App.hash);
                request.send(formData);
            });
        }
    }
</script>

{#if $App.info == null}
    Ctrl + Click on the image to add text.
{:else}
    {$App.info}
{/if}
<div class="input">
    <input
        id="value"
        type="text"
        value={$App.value}
        on:keyup={handleValueChange}
    />
</div>

B:<input
    type="checkbox"
    on:click={(e) => {
        App.SetBold(e.target.checked);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }}
/>
I:<input
    type="checkbox"
    on:click={(e) => {
        App.SetItalic(e.target.checked);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }}
/>
S:<input
    type="checkbox"
    on:click={(e) => {
        App.SetSmallCaps(e.target.checked);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }}
/>
<br />
<select on:change={handleSelectText}>
    {#each $App.texts as text}
        <option value={text.id}>
            {text.value}
        </option>
    {/each}
</select>

<select on:change={handleSelectFont}>
    {#each $App.fonts as fontitem}
        <option value={fontitem}>
            {fontitem}
        </option>
    {/each}
</select>

<select on:change={handleUpdateTextAlign}>
    <option value="center"> Center </option>
    <option value="left"> Left </option>
    <option value="right"> Right </option>
</select>

<div class="colorbox">
    <b>Fill Color</b>
    <div class="colors">
        {#each $App.webcolors as color}
            <p
                style="background: {color}"
                title={color}
                on:click={(e) => handleSelectColor(e, color)}
                on:contextmenu={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    App.UpdateShadow(color);
                    Canvas.Redraw($App.cvDynamic, $App.texts, true);
                }}
            />
        {/each}
    </div>
</div>

<div class="slider">
    Font Size:
    <input
        type="range"
        min="5"
        max="500"
        step="1"
        value="30"
        title="Font size"
        on:change={handleUpdateSize}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 10;
            } else {
                e.target.value -= 10;
            }
            App.UpdateSize(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<div class="slider">
    Width:
    <input
        type="range"
        min="-1"
        max="100"
        step="1"
        value="-1"
        title="Text fit width"
        on:change={handleUpdateFitwidth}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 5;
            } else {
                e.target.value -= 5;
            }
            App.UpdateFitwidth(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<div class="slider">
    Rotation:
    <input
        type="range"
        min="-90"
        max="90"
        step="5"
        value="0"
        title="Text rotation"
        on:change={(e) => {
            App.SetAngle(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 5;
            } else {
                e.target.value -= 5;
            }
            App.SetAngle(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<div class="colorbox">
    <b>Outline Color</b>
    <div class="colors">
        {#each $App.webcolors as color}
            <p
                style="background: {color}"
                title={color}
                on:click={(e) => {
                    App.UpdateShadow(color);
                    Canvas.Redraw($App.cvDynamic, $App.texts, true);
                }}
            />
        {/each}
    </div>
</div>

<div class="slider">
    Stroke:
    <input
        type="range"
        min="0"
        max="20"
        step="1"
        value="0"
        title="Stroke thickness"
        on:change={(e) => {
            App.UpdateThickness(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 1;
            } else {
                e.target.value -= 1;
            }
            App.UpdateThickness(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<div class="slider">
    Shadow:
    <input
        type="range"
        min="0"
        max="100"
        step="1"
        value="0"
        title="Blur"
        on:change={(e) => {
            App.UpdateBlur(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 5;
            } else {
                e.target.value -= 5;
            }
            App.UpdateBlur(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<div class="slider">
    Shadow X:
    <input
        type="range"
        min="-20"
        max="20"
        step="1"
        value="0"
        title="Shadow offset X"
        on:change={(e) => {
            App.UpdateShadowX(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 1;
            } else {
                e.target.value -= 1;
            }
            App.UpdateShadowX(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<div class="slider">
    Shadow Y:
    <input
        type="range"
        min="-20"
        max="20"
        step="1"
        value="0"
        title="Shadow offset Y"
        on:change={(e) => {
            App.UpdateShadowY(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
        on:mousewheel={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (e.deltaY < 0) {
                e.target.valueAsNumber += 1;
            } else {
                e.target.value -= 1;
            }
            App.UpdateShadowY(e.target.value);
            Canvas.Redraw($App.cvDynamic, $App.texts, true);
        }}
    />
</div>

<input type="button" class="btn" value="Save" on:click={handleSaveClick} />

<style>
    #value,
    select {
        width: 90%;
        font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
    }

    #value {
        text-align: center;
        margin-bottom: 10px;
        font-family: monospace;
    }

    .colorbox b {
        color: white;
        text-align: center;
        font-size: 12px;
        cursor: pointer;
        transition: border-color 0.25s linear;
        display: block;
        border: rgb(100, 100, 100) solid 1px;
    }
    .colorbox b:hover {
        border: red solid 1px;
    }
    .colorbox {
        position: relative;
        margin: 5px auto;
        width: 90%;
        min-height: 40px;
        border: 2px solid black;
        overflow: auto;
        border-radius: 0;
        transition: border-color 0.25s linear;
    }
    .colorbox:hover {
        border: rgb(165, 164, 164) solid 1px;
    }
    .colors {
        position: relative;
        display: flex;
        flex-flow: row wrap;
        justify-content: center;
        align-items: flex-start;
    }
    p {
        float: left;
        width: 10px;
        height: 10px;
        margin: 2px;
        padding: 2px;
    }
    p:hover {
        cursor: default;
        border: 1px solid rgb(252, 37, 37);
    }
    p:active {
        transform: translateY(2px);
    }

    input[type="range"] {
        display: block;
        margin: 0;
        padding: 0;
        font-size: inherit;
        width: 100%;
        height: 1em;
        background-color: #242424;
        overflow: hidden;
        transition: border 0.3s linear;
        border-radius: 0.25em;
        border: 1px solid #242424;
    }
    .slider {
        padding: 5px 1em;
        text-align: left;
    }

    .btn {
        display: inline-block;
        background: #000;
        color: #fff;
        border: none;
        padding: 10px 20px;
        margin: 5px;
        border-radius: 5px;
        cursor: pointer;
        text-decoration: none;
        font-size: 15px;
        font-family: inherit;
    }
    .btn:hover {
        background-color: rgb(139, 139, 139);
    }
    .btn:focus {
        outline: none;
    }
    .btn:active {
        transform: scale(0.9);
    }
</style>
