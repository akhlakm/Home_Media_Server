export const Canvas = (function() {
    // wrap a text to a certain width
    // return the texts final width and height
    function printAtWordWrap(context, text, x, y, lineHeight, fitWidth, stroke) {
        fitWidth = fitWidth || 0;
        stroke = stroke || false;

        // no max text width was given
        if (fitWidth <= 0) {
            if (stroke) context.strokeText(text, x, y);
            else context.fillText(text, x, y);
            let size = context.measureText(text);
            return [size.width, lineHeight];
        }
        // to fit the max width, we will write word by word
        var words = text.split(" ");
        var currentLine = 0;
        var idx = 1;
        var width = 0;
        var maxwidth = 0;
        // for each word in the current list
        while (words.length > 0 && idx <= words.length) {
            // see how many words we can fit,
            // start with just a single word, see if it fits.
            var str = words.slice(0, idx).join(" ");
            width = context.measureText(str).width;

            if (width > maxwidth) maxwidth = width;

            // words dont fit anymore, lets write the current ones that fit
            if (width > fitWidth) {
                // we will always write at least one word on the current line.
                // increase by 1 for slicing a single word.
                if (idx == 1) {
                    idx = 2;
                }

                // Create and write the line without the one word that does not fit.
                var linetowrite = words.slice(0, idx - 1).join(" ");
                if (stroke) {
                    context.strokeText(linetowrite, x, y + lineHeight * currentLine);
                } else context.fillText(linetowrite, x, y + lineHeight * currentLine);

                // we will move on with the remaining words in the next line.
                currentLine++;
                // if we are going beyond the canvas area, lets skip writing rest of the words.
                if (y + lineHeight * currentLine > context.canvas.height) {
                    return [maxwidth, lineHeight * currentLine];
                }
                // written the current line, remove the current line words
                words = words.splice(idx - 1);
                // reset the word counter
                idx = 1;
            } else {
                // if the words fit, try with one more word ...
                idx++;
            }
        }
        if (idx > 0) {
            var linetowrite = words.join(" ");
            if (stroke)
                context.strokeText(linetowrite, x, y + lineHeight * currentLine);
            else context.fillText(linetowrite, x, y + lineHeight * currentLine);
        }

        return [maxwidth, lineHeight * (currentLine + 1)];
    }

    /** Color code to rgb list. */
    function hex2rgb(hex) {
        return [
            ("0x" + hex[1] + hex[2]) | 0,
            ("0x" + hex[3] + hex[4]) | 0,
            ("0x" + hex[5] + hex[6]) | 0,
        ];
    }

    function textHit(coord, text) {
        var buffer = text.height / 2;
        if (text.align === "left") {
            return (
                // l, r, t, b
                (coord.x > text.x - buffer) &&
                (coord.x < text.x + text.width + buffer) &&
                (coord.y > text.y - text.height) &&
                (coord.y < text.y + buffer)
            );
        } else if (text.align === "right") {
            return (
                // l, r, t, b
                (coord.x > text.x - text.width - buffer) &&
                (coord.x < text.x + buffer) &&
                (coord.y > text.y - text.height) &&
                (coord.y < text.y + buffer)
            );
        } else {
            return (
                // l, r, t, b
                (coord.x > text.x - text.width / 2 - buffer) &&
                (coord.x < text.x + text.width / 2 + buffer) &&
                (coord.y > text.y - text.height) &&
                (coord.y < text.y + buffer)
            );
        }
    }

    return { // public interface
        GetTextHit: function(textlist, coord) {
            for (var i = 0; i < textlist.length; i++) {
                if (textlist[i].value.length > 0) {
                    if (textHit(coord, textlist[i])) {
                        return i;
                    }
                }
            }

            return -1;
        },

        GetNewText: function(textList, x, y) {
            let newtext = {}
            newtext = {
                id: -1,
                value: "",
                x: 20,
                y: (textList.length + 1) * 50,
                color: "red",
                shadow: "black",
                blur: 0,
                thickness: 0,
                size: 30,
                font: "sans",
                width: 0,
                height: 0,
                fitwidth: -1,
                align: "center",
                alpha: 255,
                bold: false,
                smallcaps: false,
                italic: false,
            }
            newtext.id = textList.length;
            newtext.value = "text" + textList.length;
            newtext.x = x;
            newtext.y = y;
            if (textList.length > 0) {
                var oldtext = textList[textList.length - 1];
                newtext.color = oldtext.color;
                newtext.shadow = oldtext.shadow;
                newtext.blur = oldtext.blur;
                newtext.thickness = oldtext.thickness;
                newtext.size = oldtext.size;
                newtext.font = oldtext.font;
            }
            return newtext;
        },

        // Load the image
        LoadImage: function(imagePath, canvasid, cb = null) {
            if (imagePath == null) return;
            var canvas = document.getElementById(canvasid),
                ctx = canvas.getContext("2d");
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            var img = new Image();
            // force reload
            var r = Math.random().toString(36).substring(2, 5);
            img.src = imagePath + "?" + r;
            img.crossOrigin = "anonymous";
            img.onload = function() {
                console.log("Image:", this.width, "x", this.height, "px");
                canvas.width = canvas.parentElement.clientWidth;
                var imageRatio = this.height / this.width;
                canvas.height = imageRatio * canvas.width;
                ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
                // add opacity
                ctx.fillStyle = "rgba(0,0,0,0.1)";
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                if (cb !== null) cb(img);
            };
        },

        GetSize: function(canvasid) {
            var canvas = document.getElementById(canvasid);
            return { width: canvas.width, height: canvas.height };
        },

        SetSize: function(canvasid, dimensions) {
            var canvas = document.getElementById(canvasid);
            canvas.width = dimensions.width;
            canvas.height = dimensions.height;
        },

        GetBlob: function(canvasid, cb) {
            var canvas = document.getElementById(canvasid);
            return canvas.toBlob(cb);
        },

        // redraw a canvas with the text lines
        Redraw: function(canvasid, textList, clear = true) {
            console.log(textList);
            var canvas = document.getElementById(canvasid);
            var ctx = canvas.getContext("2d");
            if (clear) ctx.clearRect(0, 0, canvas.width, canvas.height);

            for (var i = 0; i < textList.length; i++) {
                var text = textList[i];
                if (text.value.length === 0) continue;
                var fontstyle = text.size + "px " + text.font;
                if (text.bold) fontstyle = "bold " + fontstyle;
                if (text.smallcaps) fontstyle = "small-caps " + fontstyle;
                if (text.italic) fontstyle = "italic " + fontstyle;
                console.log("fontstyle:", fontstyle)
                ctx.font = fontstyle;
                ctx.textAlign = text.align;

                if (text.thickness > 0) {
                    ctx.strokeStyle = text.shadow;
                    ctx.lineWidth = text.thickness;
                }
                if (text.blur > 0) {
                    ctx.shadowColor = text.shadow;
                    ctx.shadowBlur = text.blur;
                }
                if (text.blur > 0 || text.thickness > 0)
                    printAtWordWrap(
                        ctx,
                        text.value,
                        text.x,
                        text.y,
                        text.size,
                        (canvas.width * text.fitwidth) / 100,
                        true
                    );
                ctx.shadowBlur = 0;

                var rgb = hex2rgb(text.color);
                var alpha = text.alpha / 255;
                ctx.fillStyle = `rgba(${rgb[0]}, ${rgb[1]}, ${rgb[2]}, ${alpha})`;
                console.log("RGB:", rgb, "HEX:", text.color, "ALPHA:", alpha);
                console.log("fillStyle:", ctx.fillStyle);
                var dim = printAtWordWrap(
                    ctx,
                    text.value,
                    text.x,
                    text.y,
                    text.size,
                    (canvas.width * text.fitwidth) / 100,
                    false
                );
                textList[i].width = dim[0];
                textList[i].height = dim[1];
            }

            return textList;
        },
    };
})();