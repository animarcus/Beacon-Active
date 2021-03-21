<template>
  <q-page class="container">
    <div class="row">
      <div v-if="!showCamera" class="col-12 text-center q-pt-md">
        <img alt="Quasar logo" src="/qr_code.svg" style="width: 340px" />
      </div>
    </div>
    <div class="row justify-center q-pt-lg">
      <div class="col-12 text-center">
        <span class="text-subtitle2 text-grey-9">
          {{ textInfo }}
        </span>

        <q-btn
          color="blue-grey-10"
          rounded
          icon="camera_alt"
          label="Read QRCode"
          class="full-width"
          size="lg"
          @click="turnCameraOn()"
          v-show="!showCamera"
        />

        <p class="text-subtitle1" v-if="result">
          Last result: <b>{{ result.user_id }}</b>
        </p>
        <div
          @click="sendToBeacon"
          v-if="showCamera"
          class="q-pa-sm column">
          <qrcode-vue
            :value="text"
            :size="size"
            level="H"
            class="qrcode-view" />
          <qrcode-stream :camera="camera" @decode="onDecode" style="width:300px" class="self-center q-pb-sm">  </qrcode-stream>

      </div>
      </div>
    </div>
    <div class="row justify-center">
      <div class="col-12 q-pt-md text-center"></div>
    </div>
  </q-page>
</template>

<script>
import { QrcodeStream } from "vue-qrcode-reader";
import QrcodeVue from "qrcode.vue";
import { api } from 'boot/axios'

export default {
  name: "Generator",
  data() {
    return {
      text: '{"user_id":42}',
      size: 300,
      isValid: undefined,
      camera: "auto",
      result: undefined,
      showCamera: false
    };
  },
  components: {
    QrcodeVue,
    QrcodeStream
  },
  computed: {
    textInfo() {
      return this.showCamera
        ? "position the qrcode on the camera"
        : "Press the button and scan a qrcode.";
    }
  },
  methods: {
    async onInit(promise) {
      // show loading indicator
      try {
        const { capabilities } = await promise;
        // successfully initialized
      } catch (error) {
        if (error.name === "NotAllowedError") {
          // user denied camera access permisson
        } else if (error.name === "NotFoundError") {
          // no suitable camera device installed
        } else if (error.name === "NotSupportedError") {
          // page is not served over HTTPS (or localhost)
        } else if (error.name === "NotReadableError") {
          // maybe camera is already in use
        } else if (error.name === "OverconstrainedError") {
          // did you requested the front camera although there is none?
        } else if (error.name === "StreamApiNotSupportedError") {
          // browser seems to be lacking features
        }
      } finally {
        this.showCamera = true
        // this.camera = "auto";
        // hide loading indicator
      }
    },
    onDecode(content) {
      this.result = JSON.parse(content);
      
      api.post('http://philips-macbook.local:8080/checkin', content)
    },
    turnCameraOn() {
      this.showCamera = true;
    },
    turnCameraOff() {
      // this.camera = "off";
      this.showCamera = false;
    },
    sendToBeacon() {
      const content = {'user_id': 42}
      api.post('http://adrian-macbook.local:8080/checkin', content)
    }
  },
  created() {
    this.onInit()
  }
};
</script>
<style lang="sass">
.camera
  width: 300px
  justify-self: center
</style>