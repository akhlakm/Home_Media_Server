import { writable } from "svelte/store";

const addStylesheetURL = (url) => {
    var link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = url;
    document.getElementsByTagName("head")[0].appendChild(link);
};

const GoogleFonts = () => {
    var fontlist = [
        "Anton",
        "Secular One",
        "Viga",
        "Jua",
        "Tangerine",
        "Courgette",
        "Dokdo",
        "Suez One",
        "Quando",
        "Coustard",
    ];
    addStylesheetURL(
        "https://fonts.googleapis.com/css2?family=" + fontlist.join("&family=")
    );
    return fontlist;
};
const webfonts = () => {
    var fonts = [
        "Arial",
        "Verdana",
        "Helvetica",
        "Tahoma",
        "Trebuchet MS",
        "Times New Roman",
        "Georgia",
        "Garamond",
        "Courier New",
        "Brush Script MT",
    ];
    return fonts.concat(GoogleFonts());
};

const AppModel = function() {
    let data = {
        hash: null,
        url: null,
        texts: [],
        selected: null,
        value: null,
        cvStatic: "static-canvas",
        cvDynamic: "dynamic-canvas",
        info: null,
        webcolors: [
            "#000000", "#00CED1", "#0000FF", "#00FFFF",
            "#5C4033", "#7A5299", "#7C7C40", "#00008B",
            "#8B0000", "#008B8B", "#8B008B", "#30D5C8",
            "#90EE90", "#006400", "#008000", "#9400D3",
            "#404040", "#808080", "#808080", "#976638",
            "#996600", "#A52A2A", "#ADD8E6", "#AFAFAF",
            "#AFE4DE", "#B8860B", "#C0C0C0", "#D3D3D3",
            "#D9A465", "#DAA520", "#E0FFFF", "#E1E1E1",
            "#E75480", "#ECD9B0", "#EE82EE", "#EEBC1D",
            "#EEDD62", "#F0DC82", "#F1E5AC", "#F2E58F",
            "#FF0000", "#FF00FF", "#FF8C00", "#FF77FF",
            "#FF3333", "#FFA500", "#FFB6C1", "#FFC0CB",
            "#FFCC00", "#FFD700", "#FFDB58", "#FFEC8B",
            "#FFF8C9", "#FFFF00", "#FFFFE0", "#FFFFF0", "#FFFFFF",
        ],
        fonts: webfonts(),
    }
    const { subscribe, set, update } = writable(data);
    var model = {
        subscribe,
        ImageUrl: url => update(d => { d.url = url; return d; }),
        ImageHash: hash => update(d => { d.hash = hash; return d; }),
        SelectText: i => update(d => {
            d.selected = i;
            d.value = d.texts[d.selected].value;
            return d;
        }),
        UpdateText: t => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].value = t;
            d.value = t;
            return d;
        }),
        UpdatePosition: (dx, dy) => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].x += dx;
            d.texts[d.selected].y += dy;
            return d;
        }),
        UpdateFont: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].font = f;
            return d;
        }),
        SetBold: b => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].bold = b;
            return d;
        }),
        SetItalic: b => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].italic = b;
            return d;
        }),
        SetSmallCaps: b => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].smallcaps = b;
            return d;
        }),
        UpdateSize: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].size = f;
            return d;
        }),
        UpdateThickness: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].thickness = f;
            return d;
        }),
        UpdateBlur: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].blur = f;
            return d;
        }),
        UpdateAlign: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].align = f;
            return d;
        }),
        UpdateShadow: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].shadow = f;
            return d;
        }),
        UpdateFitwidth: f => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].fitwidth = f;
            return d;
        }),
        SelectColor: c => update(d => {
            if (d.selected === null) return d;
            d.texts[d.selected].color = c;
            return d;
        }),
        AddText: t => update(d => {
            d.texts = [...d.texts, t];
            d.selected = d.texts.length - 1;
            d.value = t.value;
            return d;
        }),
        SetInfo: m => update(d => {
            d.info = m;
            setTimeout(() => {
                model.SetInfo(null);
            }, 3000);
            return d;
        }),
    }

    return model;
}();

export default AppModel;