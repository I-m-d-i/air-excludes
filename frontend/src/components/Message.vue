<template>
  <div style="    top: 0;
    pointer-events: none;
    z-index: 1000;
    width: 100%;
    height: 100vh;
    position: absolute;">
    <v-snackbars :key="key" :objects.sync="objects">
      <template v-slot:default="{ message, id, index }">
          <span style="max-height:600px; white-space: pre-line; font-size: 0.85rem;">
         {{ message }}
          </span>
      </template>
      <template v-slot:action="{ close }">
        <v-icon @click="close">
          mdi-close-circle-outline
        </v-icon>
      </template>
    </v-snackbars>
  </div>
</template>
<script>
import VSnackbars from "v-snackbars";
import {bus} from "@/main";

export default {
  components: {
    "v-snackbars": VSnackbars,
  },
  name: "messages",
  data: () => ({objects: [], key: Date.now() + (Math.random() + "").slice(2)}),
  methods: {
    addMessage(message) {
      this.objects.push({
        message: message.message,
        top: message.top !== undefined,
        bottom: message.bottom === undefined ? true : message.bottom,
        left: message.left !== undefined,
        right: message.right !== undefined,
        color: message.color || "blue",
        transition: message.transition || "scale-transition",
        timeout: message.timeout === undefined ? 5000 : message.timeout,
      });
    },
    reset(){
      this.objects = []
      this.key = Date.now() + (Math.random() + "").slice(2)
    }
  },
  mounted() {
    bus.$on('message-add', async (message) => {
      this.$nextTick(() => {
        this.addMessage(message)
      });
    })
    bus.$on('message', (message) => {
      if (this.objects.length > 0) {
        this.reset()
      }
      this.$nextTick(()=>this.addMessage(message))
    })
    bus.$on('message-reset', () => {
      this.reset()
    })
  },
}
</script>
<!--
<template>
  <div style="
    top: 0;
    pointer-events: none;
    z-index: 1000;
    display: flex;
    justify-content: center;
    width: 100%;
    height: 100vh;
    position: absolute;
    flex-direction: column;
justify-content: flex-end;"
  >
    <v-snackbar
        id="bar"
        v-for="(item,index) in listMessage"
        :key="listMessage.length"
        max-height="600"
        :centered="item.centered"
        multi-line
        :timeout="-1"
        v-model="item.show"
        :color="item.color"
        transition="scale-transition"
    >
          <span v-bind:style="{color:item.textColor }"
                style="max-height:600px; white-space: pre-line; font-size: 0.85rem;">
         {{ item.message }}
          </span>
      <template v-slot:action="{ attrs }">
        <v-icon
            :color="item.textColor"
            v-bind="attrs"
            @click="deleteMessage(item)">
          mdi-close-circle-outline
        </v-icon>
      </template>
    </v-snackbar>
  </div>
</template>
<script>
import {bus} from "@/main";

export default {
  name: "message",
  data: () => ({
    active: 0,
    listMessage: [],
    centered: false,
    left: false,
    right: false,
    top: false,
    bottom: false,
    stack: false,
    messagesDisplay: {
      message: "",
      color: null,
      show: false,
      textColor: null,
      timeout: 0,

    },
    timeoutID: null,
  }),
  mounted() {
    bus.$on('message', (message) => {
      let messagesDisplay = {
        message: "",
        color: null,
        show: false,
        textColor: null,
        centered: false,
        timeout: 0,
        timeoutID: null
      }

      if (!(message.show === undefined || message.show)) {
        if (this.timeoutID != null) {
          clearTimeout(this.timeoutID)
          this.close()
        }
        return;
      }
      messagesDisplay.show = (message.show === undefined || message.show)
      messagesDisplay.message = (message.message === undefined ? "" : message.message)
      messagesDisplay.color = message.color || "#4caf50"
      messagesDisplay.textColor = message.textColor || "FFFFFFFF"
      messagesDisplay.centered = message.centered || false
      messagesDisplay.timeout = message.timeout || 4000
      //Если установлен режим stack, то добавляем в стэк сооющения, иначе сбрасываем список и добавляем новое сообщение
      if (this.stack) {
        this.listMessage.push(messagesDisplay)
      } else {
        if (this.listMessage.length > 0) {
          this.listMessage = []
          messagesDisplay.show = false
          setTimeout(() => {
            messagesDisplay.show = true
          }, 200)
          clearTimeout(messagesDisplay.timeoutID)
        }
        this.listMessage.push(messagesDisplay)
      }
      //Заводим таймер если надо (если -1 то не ставим таймер)
      if (messagesDisplay.timeout !== -1) {
        messagesDisplay.timeoutID = setTimeout(() => {
          this.close(messagesDisplay)
        }, messagesDisplay.timeout)
      }
    })
    bus.$on('message-option', (options) => {
      this.stack = options.stack === undefined || options.stack
      this.right = options.right === undefined || options.right
      this.centered = options.center === undefined || options.center
      this.left = options.left === undefined || options.left
      this.top = options.top === undefined || options.top
      this.bottom = options.bottom === undefined || options.bottom
    })
  },
  methods: {
    close(message) {
      message.show = false
      setTimeout(() => {
        this.deleteMessage(message)
      }, 200)
    },
    deleteMessage(message) {
      for (let i = 0; i < this.listMessage.length; i++) {
        if (this.listMessage[i] === message) {
          this.listMessage.splice(i, 1)
        }
      }
    }
  },
}
</script>
<style scoped>
.v-snack {
  position: relative;
  height: max-content;
}

#bar {

}
</style>


-->
