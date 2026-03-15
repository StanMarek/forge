function copyOutput() {
    var output = document.querySelector('#output textarea');
    if (output) {
        navigator.clipboard.writeText(output.value);
        var btn = document.querySelector('[data-copy-btn]');
        if (btn) {
            var original = btn.textContent;
            btn.textContent = 'Copied!';
            setTimeout(function() { btn.textContent = original; }, 1500);
        }
    }
}
