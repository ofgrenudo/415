function move_poster() {
  const poster = document.getElementById('poster');
  const header = document.querySelector('header');
  const footer = document.querySelector('footer');
  possible_posters = [
    "/img/propaganda/build_communities.jpg",
    "/img/propaganda/care_not_control.jpg",
    "/img/propaganda/individualism_destroys.jpg",
    "/img/propaganda/nobosses.jpg",
    "/img/propaganda/palestine_will_be_free.jpg",
    "/img/propaganda/support_indigenous_women.jpg",
    "/img/propaganda/liberation_belongs_to_us.jpg",
    "/img/propaganda/wage_theft.jpg",
    "/img/propaganda/workers_run_the_world.jpg",
    //"/img/propaganda/american_centry_of_humilitation.jpg",
    "/img/propaganda/american_commies.jpg",
    "/img/propaganda/american_infidel.jpg",
    "/img/propaganda/bamna.jpg",
    "/img/propaganda/born_just_in_time.jpg",
    "/img/propaganda/every_day_i_dream_about_b2b_software.jpg",
    "/img/propaganda/free_busses.jpg",
    "/img/propaganda/if_g_made_the_iphone.jpg",
    "/img/propaganda/mamdani_our_commie.jpg",
    //"/img/propaganda/mamdani_unlimited_genocide.jpg",
    "/img/propaganda/year_of_the_commies.jpg"
  ]
  const magic_number = Math.floor(Math.random() * possible_posters.length);

  console.log("Selected propaganda poster of the load is " + possible_posters[magic_number]);
  poster.src = possible_posters[magic_number];

  if (!poster || !header || !footer) return;

  const isMobile = window.matchMedia('(max-width: 767px)').matches;

  if (isMobile) {
    if (poster.parentNode !== footer) {
      footer.insertBefore(poster, footer.firstChild);
    }
    // add top spacing on mobile
    poster.style.marginTop = '12px'; // or '1rem'
  } else {
    const hrs = header.querySelectorAll('hr');
    const targetHr = hrs[hrs.length - 1];
    if (!targetHr) return;

    if (poster.parentNode !== header || poster.previousElementSibling !== targetHr) {
      targetHr.parentNode.insertBefore(poster, targetHr.nextSibling);
    }
    // remove the mobile spacing when not mobile
    poster.style.marginTop = '';
  }
}

function pick_a_header_gif() {
    console.log()
    var header_image = document.getElementById("header_gif");
    possible_headers = ["/gifs/mega_man.gif", "/gifs/welcome_home.gif", "/gifs/powerball.gif", "/gifs/xxx.gif", "/gifs/bookmark_me.gif", "/gifs/butterflys.gif", "/gifs/got_linux.gif"]
    const magic_number = Math.floor(Math.random() * possible_headers.length);

    console.log("Selected header of the load is " + possible_headers[magic_number]);
    header_image.src = possible_headers[magic_number];

}

function resize_todo_list() {
    var todo_iframe = document.getElementById("shit_i_need_to_do");
    var papa_todo = document.getElementById("papa_todo_list");
    if (todo_iframe != null) {
        console.log("Setting todo iframe to: " + todo_iframe.contentWindow.document.documentElement.scrollHeight);
        todo_iframe.style.height = todo_iframe.contentWindow.document.documentElement.scrollHeight + 'px';
        if (todo_iframe.contentWindow.document.documentElement.scrollHeight == 0) {
            console.log("Hiding papa_todo");
            papa_todo.hidden = true
        }
    }
}

function resize_change_log() {
    var changelog_iframe = document.getElementById("shit_i_did");
    var papa_change_log = document.getElementById("papa_change_log");
    if (changelog_iframe != null) {
        console.log("Setting change log iframe to: " + changelog_iframe.contentWindow.document.documentElement.scrollHeight);
        changelog_iframe.style.height = changelog_iframe.contentWindow.document.documentElement.scrollHeight + 'px';
        if (changelog_iframe.contentWindow.document.documentElement.scrollHeight == 0) {
            console.log("Hiding papa_change_log")
            papa_change_log.hidden = true
        }
    }
}

function onload_scheduler(){
    console.log("picking a header");
    pick_a_header_gif();
    console.log("setting todo size");
    resize_todo_list();
    console.log("setting change log size");
    resize_change_log();
    console.log("moving poster");
    move_poster();
}

window.onload = onload_scheduler;
