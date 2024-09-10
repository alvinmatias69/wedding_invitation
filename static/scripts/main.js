const HIDE_CHURCH_PARAMS = "hc";
const CHURCH_ID = "church-events";

// Add your javascript here
// Don't forget to add it into respective layouts where this js file is needed
$(document).ready(function() {
  $('#liveToast').toast({});
  $('#notifToast').toast({});

  AOS.init({
    // uncomment below for on-scroll animations to played only once
    once: true
  }); // initialize animate on scroll library

  const urlParams = new URLSearchParams(window.location.search)
  if (urlParams.has(HIDE_CHURCH_PARAMS) && urlParams.get(HIDE_CHURCH_PARAMS) === "y") {
    document.getElementById(CHURCH_ID).remove();
  }

  document.getElementById("donation-bca").addEventListener("click", () => handle_donation("7751045852"));
  document.getElementById("donation-ocbc").addEventListener("click", () => handle_donation("015810061463"));
});

// Smooth scroll for links with hashes
$("a.smooth-scroll").click(function(event) {
  // On-page links
  if (
    location.pathname.replace(/^\//, "") == this.pathname.replace(/^\//, "") &&
    location.hostname == this.hostname
  ) {
    // Figure out element to scroll to
    var target = $(this.hash);
    target = target.length ? target : $("[name=" + this.hash.slice(1) + "]");
    // Does a scroll target exist?
    if (target.length) {
      // Only prevent default if animation is actually gonna happen
      event.preventDefault();
      $("html, body").animate(
        {
          scrollTop: target.offset().top
        },
        1000,
        function() {
          // Callback after animation
          // Must change focus!
          var $target = $(target);
          $target.focus();
          if ($target.is(":focus")) {
            // Checking if the target was focused
            return false;
          } else {
            $target.attr("tabindex", "-1"); // Adding tabindex for elements not focusable
            $target.focus(); // Set focus again
          }
        }
      );
    }
  }
});

// Photo Filter
var activeFilter = "all";

$(".ww-filter-button").on("click", function(e) {
  // remove btn-primary from all buttons first
  $(".ww-filter-button").removeClass("btn-primary");
  $(".ww-filter-button").addClass("btn-outline-primary");

  // add btn-primary to active button
  var button = $(this);
  button.removeClass("btn-outline-primary");
  button.addClass("btn-primary");
  filterItems(button.data("filter"));
  e.preventDefault();
});

function filterItems(filter) {
  if (filter === activeFilter) {
    return;
  }

  activeFilter = filter;
  $(".ww-gallery .card").each(function() {
    var card = $(this);
    var groups = card.data("groups");
    var show = false;
    if (filter === "all") {
      show = true;
    } else {
      for (var i = 0; i < groups.length; i++) {
        if (groups[i] === filter) {
          show = true;
        }
      }
    }
    // hide everything first
    card.fadeOut(400);
    setTimeout(function() {
      if (show && !card.is(":visible")) {
        card.fadeIn(400);
      }
    }, 500);
  });
}

// Light Box
$(document).on("click", '[data-toggle="lightbox"]', function(event) {
  event.preventDefault();
  $(this).ekkoLightbox();
});

function handle_donation(payload) {
  const textArea = document.createElement('textarea');
  textArea.value = payload;
  textArea.style.opacity = 0;
  document.body.appendChild(textArea);
  textArea.select();
  try {
    document.execCommand('copy');
    $('#liveToast').toast('show');
  } catch (err) {
    console.error(err.name, err.message);
  } finally {
    document.body.removeChild(textArea);
  }
}
