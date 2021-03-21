<template>
  <q-page class="q-px-lg">
    <h5 class="fontsize-25">What is this website anyway?</h5>
    <p class="fontsize-12 bg-secondary "> This project is the result of the 2021 36-hours START-hack hackathon. <br>
        A team of three brothers created this for the Helsana case aiming to make Switzerland fitter.</p>
    <p class="fontsize-10"> <span class="textpointers fontsize-14">What is the aim?</span> <br> 
        Our project makes sure people go out for a walk during the day while working from home.
        To do this, we need incentive!
        By placing "beacons" around the city containing a QR code scanner and a display, we can make people log their time outside and be placed on a leaderboard that can be created by companies, campuses and more.</p>

    <p class="fontsize-10"> <span class="textpointers fontsize-14">How does it work?</span> <br> 
        On the first launch of the app, the user must input a user name, and a unique identifier is automatically assigned.
        These details are then sent to the server and stored safely.
        When a user goes to one of these beacons, they must scan the QR code appearing on screen.<br>
        The beacon processes the request, and issues a cryptographic token authenticating the checkin at the current time.
        This token appears on the beacon's screen, and the user must scan it to relay it to a server along with the user's ID and timestamps.<br>
        With few assumptions on the beacon, the server is convinced that this token is valid (via digital signatures).
        If the user performs both a check in and a check out, the server can also deduce the amount of this user spent outside.</p>

    <p class="fontsize-10"> <span class="textpointers fontsize-14">How does it scale?</span> <br> 
        Since this project requires physical deployment and maintenance of beacons, we envision these beacons to be be cheap to make (using a Raspberry Pi with minimal accessories).
        More importantly, they do not need to be connected to the internet.
        Some fraud prevention mechanism can also be easily be implemented by monitoring how this scanning is done like using on-device machine learning to scan faces, and logging for a short period of time the users accessing the beacons.<br>
        The server itself can also perform checks, and only accept tokens within reasonable timestamps (eg user can check in only once per day). <br>
        Finally, these beacons can be operated independently from the app and the server.
        For example, a local city could decide to deploy and maintain the beacons, and license other companies to use them.</p>

    <p class="fontsize-10"> <span class="textpointers fontsize-14">What's more to come?</span> <br> 
        This is the base idea of the project but to motivate people even more, we want to add the possibility to create leaderboards that users can signup for.
        If a school wants to motivate it's students to go outside during their long hours studying, it can place these beacons around campus and have a leaderboard so that students can motivate eachother with friendly competition for how much time spent outside.</p>

    <p class="fontsize-10"> <span class="textpointers fontsize-14">How was the project created?</span> <br> 
        We're using Javascript along with the Quasar framework for maximum accessibility across devices. The server and the beacon are programmed with Go to handle the cryptography side of things. 
            <br><br><br>
        <span class="fontsize-10">Our goal with this project is to enable healthier working from home / at school by a simple gesture which is taking a walk outside.
        Other apps that log physical activity can be demotivating because some people are way more active than others. With our simple system of just logging time outside, people can be motivated more easily and companies / schools can provide rewards for the users.</span></p>
  </q-page>
</template>

<script>
import { LocalStorage } from 'quasar'



export default {
  name: 'Info',
  // components: { EssentialLink },
  data () {
    return {
      firstTime: this.getStorage()
    }
  },
  methods: {
    getStorage() {
      return LocalStorage.isEmpty()
    },
    redirect() {
        console.log("aaaaaaaaaaaaaaa")
        if (!this.firstTime) return
        const pk = this.generatePk()

        LocalStorage.set('pk', pk)
        this.$router.push({path: '/signup'})
    },
    generatePk() {
        return 'aaaaaaa'
    }
  },
  created () {
    this.redirect()
  }
}
</script>


<style lang="scss">
    .textpointers {
        font-weight: bold;

    }
</style>