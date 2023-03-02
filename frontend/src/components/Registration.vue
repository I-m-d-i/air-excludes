<template>
  <div style="display: flex; align-items: center;
    justify-content: center; padding-block: 25vh;">
    <v-card
        min-height="300"
        width="400"
        outlined>
      <v-card-title style="flex-direction: column; text-align: center;">
        <h1 style="
          justify-content: space-around;
          text-align: center;
          font-size: 1.75em;
          color: rgb(100 165 19);
          font-family: cursive;">
          Krasecology.ru
        </h1>
        <h2 style="font-weight: 400; padding-top: 16px; font-size: 24px;">
          <span>Регистрация</span>
        </h2>
      </v-card-title>
      <v-form
          ref="form"
          v-model="valid"
          @submit.prevent="validate">
        <v-card-text style="padding-bottom: 0;">
          <v-text-field
              v-model="email"
              :rules="[rules.required, rules.email]"
              label="Электронная почта"
              required
          />
          <v-text-field
              v-model="pass"
              :rules="[rules.required]"
              label="Пароль"
              required
              type="password"
          />
          <v-text-field
              v-model="secondName"
              :rules="[rules.required]"
              label="Фамилия"
              required
          />
          <v-text-field
              v-model="firstName"
              :rules="[rules.required]"
              label="Имя"
              required
          />
          <v-text-field
              v-model="lastName"
              :rules="[rules.required]"
              label="Отчество"
              required
          />
          <v-text-field
              v-model="org"
              :rules="[rules.required]"
              label="Организация"
              required
          />
        </v-card-text>
        <v-card-actions style="margin-inline: 5px;">
          <v-btn @click="returnToHome()">
            Вернуться
          </v-btn>
          <v-spacer/>
          <vue-recaptcha
              ref="recaptcha"
              size="invisible"
              sitekey="6Ld__vUiAAAAAJzOyCAGzHswCeCL6V5dIE4x-X_k"
              @verify="registration"
              @expired="onCaptchaExpired">
            <v-btn type="submit" :disabled="!valid">
              Зарегистрироваться
            </v-btn>
          </vue-recaptcha>
        </v-card-actions>
      </v-form>
    </v-card>
  </div>
</template>

<script>

import axios from "axios";
import {VueRecaptcha} from 'vue-recaptcha';
import {bus} from '@/main';

export default {
  components: {VueRecaptcha},
  name: "registration",
  data() {
    return {
      rules: {
        required: value => !!value || 'Обязательное поле.',
        email: value => {
          const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
          return pattern.test(value) || 'Неверный адрес электронной почты.'
        },
      },
      valid: false,
      email: '',
      pass: '',
      origin: '',
      secondName: '',
      firstName: '',
      lastName: '',
      org: '',
    }
  },
  methods: {
    validate() {
      this.$refs.recaptcha.execute()
    },
    onCaptchaExpired() {
      this.$refs.recaptcha.reset()
    },
    registration() {
      axios
          .post("http://account.krasecology.ru/api/account/registration", {
            organisation: this.org,
            password: this.pass,
            email: this.email,
            origin: window.location.href,
            name: {
              lastname: this.secondName,
              firstname: this.firstName,
              patronymic: this.lastName,
            },
          }, {withCredentials: false})
          .then(() => {
            bus.$emit('message', {
              message: "На вашу почту " + this.email + " было отправлено письмо для подтверждения почты. " +
                  "Регистрация вашего аккаунта будет завершена после подтверждения почты " +
                  "(прохождения по ссылке в письме)",
              color: "green",
            });
            this.$router.replace("login")
          })
          .catch(() => {
            bus.$emit('message', {
              message: "Ошибка регистрации",
              color: "red",
            });
            this.onCaptchaExpired()
          })
    },
    returnToHome() {
      this.$router.push(this.$route.params.nextUrl)
    }
  }
}

</script>
