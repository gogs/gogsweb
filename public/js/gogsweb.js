$(document).ready(function () {
    // Small devices menu.
    $('#menu-link').click(function () {
        $('#wide-nav').slideToggle('fast');
    });

    $('main').click(function () {
        if ($('#menu-link').is(':visible')) {
            $('#wide-nav').slideUp('fast');
        }
    });

    // Language option button.
    $('#lang').on('change', function () {
        document.location.href = this.options[this.selectedIndex].value;
    });

    // Render code blocks.
    $('.markdown').find('pre > code').parent().addClass('prettyprint');
    prettyPrint();

    // Encode url.
    var $doc = $('.docs-markdown');
    $doc.find('a').each(function () {
        var node = $(this);
        var link = node.attr('href');
        var index = link.indexOf('#');
        if (link.indexOf('http') === 0 && link.indexOf(window.location.hostname) === -1) {
            return;
        }
        if (index < 0 || index + 1 > link.length) {
            return;
        }
        var val = link.substring(index + 1, link.length);
        val = encodeURIComponent(decodeURIComponent(val).toLowerCase().replace(/\s+/g, '-'));
        node.attr('href', link.substring(0, index) + '#' + val);
    });

    // Set anchor.
    $doc.find('h1, h2, h3, h4, h5, h6').each(function () {
        var node = $(this);
        if (node.hasClass('ui')) {
            return;
        }
        var val = encodeURIComponent(node.text().toLowerCase().replace(/\s+/g, "-"));
        node = node.wrap('<div id="' + val + '" class="anchor-wrap" ></div>');
        node.append('<a class="anchor" href="#' + val + '"><span class="octicon octicon-link"></span></a>');
    });
});