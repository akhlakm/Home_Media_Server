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
    function handleUpdate(e) {
        App.UpdateFont(e.target.value);
        Canvas.Redraw($App.cvDynamic, $App.texts, true);
    }

    function handleSaveClick(e) {
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
</script>

<div class="input">
    <input
        id="value"
        type="text"
        value={$App.value}
        on:change={handleValueChange}
    />
</div>

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

<div class="colorbox">
    <b>Color</b>
    <div class="colors">
        {#each $App.webcolors as color}
            <p
                style="background: {color}"
                title={color}
                on:click={(e) => handleSelectColor(e, color)}
            />
        {/each}
    </div>
</div>

<div class="slider">
    <input
        type="range"
        min="5"
        max="250"
        step="1"
        title="Font size"
        on:change={handleUpdateSize}
    />
</div>

<input type="button" class="btn" value="Save" on:click={handleSaveClick} />

<style>
    #value,
    select {
        width: 90%;
        font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
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
        padding: 1em;
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
