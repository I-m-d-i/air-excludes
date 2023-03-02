<template>
  <div>
    <div class="d-flex flex-column justify-center align-center" style="height: 100vh;padding-inline: 1vw;">
      <v-card
          class="d-flex-block"
          min-height="20vh"
          min-width="240px"
          max-width="450px"
          width="100%"
          outlined>
        <v-card-title style="flex-direction: column;text-align: center;">
          <h1 style="
        justify-content: space-around;
        text-align: center;
        font-size: 1.75em;
        color: rgb(100 165 19);
        font-family: cursive;">
            Krasecology.ru
          </h1>
          <h2 style="
        font-weight: 400;
        padding-top: 16px;
        font-size: 24px;">
            <span>Вход</span>
          </h2>
        </v-card-title>
        <v-form
            ref="form"
            v-model="valid"
            method="post"
            @submit.prevent="validate">
          <v-card-text style="padding-bottom: 0;">
            <v-text-field
                v-model="name"
                :rules="[rules.required, rules.email]"
                label="Электронная почта"
                required/>
            <v-text-field
                :append-icon="showPass ? 'mdi-eye' : 'mdi-eye-off'"
                v-model="pass"
                :rules="[rules.required]"
                label="Пароль"
                required
                hide-details
                :type="showPass ? 'text' : 'password'"
                @click:append="showPass = !showPass"/>
            <a style="font-size: 9pt" href="http://account.krasecology.ru/#/password-reset">Забыли пароль?</a>
          </v-card-text>
          <v-card-actions style="margin-inline: 5px; margin-block: 5px">
            <v-spacer/>
            <vue-recaptcha
                ref="recaptcha"
                size="invisible"
                sitekey="6Ld__vUiAAAAAJzOyCAGzHswCeCL6V5dIE4x-X_k"
                @verify="login"
                @expired="onCaptchaExpired"
            />
            <v-btn type="submit" :disabled="!valid">
              Войти
            </v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
      <div style="text-align: center; padding: 10px; color:#232629; font-size: 11pt;">
        Ещё нет учётной записи?
        <a @click="registration">Зарегистрируйтесь</a>
      </div>
    </div>
  </div>

</template>
<script>
import axios from "axios";
import {VueRecaptcha} from 'vue-recaptcha';
import {bus} from '@/main';

export default {
  components: {VueRecaptcha},
  name: 'Login',
  data() {
    return {
      showPass: false,
      valid: false,
      name: '',
      pass: '',
      rules: {
        required: value => !!value || 'Обязательное поле.',
        email: value => {
          const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
          return pattern.test(value) || 'Неверный адрес электронной почты.'
        },
      },
    }
  },
  methods: {
    validate() {
      this.$refs.recaptcha.execute()
    },
    onCaptchaExpired() {
      this.$refs.recaptcha.reset()
    },
    login(recaptchaToken) {
      //e.preventDefault()
      axios
          .post("/api/auth/login", {login: this.name, pass: this.pass, recaptchaToken: recaptchaToken})
          .then(() => {
            if (this.$route.params.nextUrl != null) {
              this.$router.replace(this.$route.params.nextUrl)
            } else {
              this.$router.replace('/')
            }
          })
          .catch(e => {
            if (e.response.status === 500 && e.response.data === "неверный логин/пароль") {
              bus.$emit('message', {
                message: "Неверная электронная почта или пароль",
                color: "red"
              });
            }
            this.onCaptchaExpired()
          });
    },
    logout() {
      axios
          .post("/api/auth/logout", {})
          .then(response => {
            console.log(response.status)
          })
          .catch(() => {
          });
    },
    registration() {
      this.$router.push({name: "registration", params: {nextUrl: this.$route.path}})
    },
  },
  created() {

  }
}
</script>