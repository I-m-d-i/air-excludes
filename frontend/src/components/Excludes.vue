<template>
  <div style="margin-inline: 15px;">
    <v-row class="ma-2">
      <v-btn
          small
          class="ma-2"
          @click="addRow()">
        Добавить
      </v-btn>
      <v-btn
          :disabled="rowsForDelete.length<=0"
          small
          class="ma-2"
          @click="multiDelete()">
        Удалить
      </v-btn>
      <v-btn
          small
          :disabled="rowsForDelete.length<=0"
          class="ma-2"
          @click="rowsForDelete=[]">
        Отменить выбор
      </v-btn>
      <v-spacer/>
      <v-btn
          :disabled="(!(changes.updatingRows.length > 0 || changes.addedRows.length > 0 || changes.deletedRows.length > 0))"
          small
          class="ma-2"
          @click="checkErrors()?'':windowSave=true">
        Сохранить
      </v-btn>
      <v-btn
          small
          class="ma-2"
          @click="(changes.updatingRows.length > 0 || changes.addedRows.length > 0 || changes.deletedRows.length > 0)?windowAlert=true:rollback()">
        Отменить
      </v-btn>
    </v-row>
    <v-data-table
        :loading="loadingData"
        v-model="rowsForDelete"
        :headers="headerTable"
        :items="filteredDesserts"
        item-key="id"
        :item-class="itemClasses"
        :page.sync="page"
        :items-per-page="20"
        @page-count="pageCount = $event"
        class="elevation-1"
        show-select
        hide-default-footer
        dense>
      <template v-slot:footer>
        <div style="display: flex; border-top:thin solid rgba(0,0,0,.12);">
          <v-spacer/>
          <v-btn small text @click="windowInportFile=true">
            Импорт данных стабильности
          </v-btn>
        </div>
      </template>
      <template v-slot:header.postId="{ header }">
        {{ header.text }}
        <v-menu>
          <template v-slot:activator="{ on, attrs }">
            <v-icon small
                    :color="postFilterValue ? 'primary' : ''"
                    v-bind="attrs"
                    v-on="on">
              mdi-filter
            </v-icon>
          </template>
          <v-list>
            <v-btn text @click="postFilterValue = null" class="ms-5 me-5 mb-2" small>Очистить фильтр
              <v-icon right>
                mdi-filter-remove
              </v-icon>
            </v-btn>
            <v-divider></v-divider>
            <v-list-item
                v-for="(post) in posts"
                :key="post.id"
                @click="postFilterValue = post.id"
                link>
              <v-list-item-content>
                <v-list-item-title>{{ post.name }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
      <template v-slot:header.sensorType="{header}">
        {{ header.text }}
        <v-menu :close-on-content-click="false">
          <template v-slot:activator="{ on, attrs }">
            <v-icon small
                    :color="sensorFilterValue.length>0 ? 'primary' : ''"
                    v-bind="attrs"
                    v-on="on">
              mdi-filter
            </v-icon>
          </template>
          <v-list>
            <v-btn text @click="sensorFilterValue=[]" class="ms-5 me-5 mb-2" small>Очистить фильтр
              <v-icon right>
                mdi-filter-remove
              </v-icon>
            </v-btn>
            <v-divider></v-divider>
            <v-list-item-group
                multiple
                v-model="sensorFilterValue">
              <v-list-item
                  v-for="(sensor) in sensorTypes"
                  :key="sensor.id"
                  :value="sensor.id">
                <template v-slot:default="{ active }">
                  <v-list-item-action>
                    <v-checkbox :input-value="active"></v-checkbox>
                  </v-list-item-action>
                  <v-list-item-content>
                    <v-list-item-title>{{ sensor.name }}</v-list-item-title>
                  </v-list-item-content>
                </template>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-menu>
      </template>
      <template v-slot:item.actions="{ item }">
        <v-btn style="margin: auto" x-small
               :title="`Копирование записи \nПКМ - множественное копирование`"
               @click="cloningRow(item,1)"
               @contextmenu="showMenu($event,item)">
          Копировать
        </v-btn>
        <v-menu
            v-if="item===menu.menuItem"
            v-model="menu.showMenu"
            :position-x="menu.x"
            :position-y="menu.y"
            absolute
            :close-on-content-click="false"
            offset-y>
          <v-card elevation="2">
            <v-card-title class="text-h5 mb-1">
              Количество копий:
            </v-card-title>
            <v-text-field min="1" suffix="шт." v-model:number="numberOfCopies" type="number"
                          class="ms-5 me-5"></v-text-field>
            <v-card-actions>
              <v-btn
                  outlined
                  rounded
                  text
                  @click="cloningRow(item,numberOfCopies)">
                Копировать
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-menu>
      </template>
      <template v-slot:item.postId="{isMobile, item, value, header}">
        <v-menu offset-y>
          <template v-slot:activator="{ on, attrs }">
            <div class="cell"
                 v-bind="attrs"
                 v-on="on"
                 v-bind:style="isMobile&&value===null?'min-width: 150px;':''">
              <span>{{ posts.filter(pItem => pItem.id === value)[0]?.name || " " }}</span>
              <v-spacer/>
              <v-icon size="10" v-if="isUpdated(item.id,header.value)" color="yellow">mdi-circle</v-icon>
              <v-icon size="10" v-if="value===null" color="red">mdi-circle</v-icon>
            </div>
          </template>
          <v-list>
            <v-list-item
                v-for="(post) in posts"
                :key="post.id"
                @click="updating(item,post.id,header.value)"
                link>
              <v-list-item-content>
                <v-list-item-title>{{ post.name }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
      <template v-slot:item.sensorType="{isMobile, item, value, header}">
        <v-menu offset-y>
          <template v-slot:activator="{ on, attrs }">
            <div
                class="cell"
                v-bind="attrs"
                v-on="on"
                v-bind:style="isMobile&&value===null?'min-width: 150px;':''"
            >
              <span style="object-fit: fill;">{{
                  sensorTypes.filter(sItem => sItem.id === value)[0]?.name || ' '
                }}</span>
              <v-spacer/>
              <v-icon size="10" v-if="isUpdated(item.id,header.value)" color="yellow">mdi-circle</v-icon>
              <v-icon size="10" v-else-if="value===null" color="red">mdi-circle</v-icon>
            </div>
          </template>
          <v-list>
            <v-list-item
                v-for="(sensorType) in sensorTypes"
                :key="sensorType.id"
                @click="updating(item,sensorType.id,header.value)"
                link>
              <v-list-item-content>
                <v-list-item-title>{{ sensorType.name }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
      <template v-slot:item.dateStart="{ item, value, header}">
        <v-text-field
            class="date cell"
            dense
            :required="true"
            :value="value"
            hide-details
            single-line
            :min="'2000-01-01T00:00'"
            type="datetime-local"
            @input="updating(item,$event, header.value)"
        >
          <template v-slot:append>
            <v-icon size="10"
                    v-if="!validateDate(item)||new Date(value)<new Date('2000-01-01T00:00')||value===null||value===''"
                    color="red"
                    class="error-dot">mdi-circle
            </v-icon>
            <v-icon size="10" v-else-if="isUpdated(item.id,header.value)" color="yellow" class="update-dot">mdi-circle</v-icon>
          </template>
        </v-text-field>
      </template>
      <template v-slot:item.dateEnd="{ item, value, header}">
        <v-text-field
            class="date cell"
            :value="value"
            dense
            hide-details
            single-line
            :min="item.dateStart"
            type="datetime-local"
            @input="updating(item,$event, header.value)"
        >
          <template  v-slot:append>
            <v-icon v-if="!validateDate(item)||new Date(value)<new Date('2000-01-01T00:00')||value===null||value===''"
                    size="10"
                    color="red"
                    class="error-dot">
              mdi-circle
            </v-icon>
            <v-icon v-else-if="isUpdated(item.id,header.value)"
                    size="10"
                    color="yellow"
                    class="update-dot">
              mdi-circle
            </v-icon>
          </template>
        </v-text-field>
      </template>
      <template v-slot:item.comment="{isMobile, item, value, header }">
        <v-edit-dialog
            save-text="Сохранить"
            cancel-text="Отменить"
            large
            v-bind:style="isMobile&&value===null?'min-width: 150px;':''"
            @open="editable=value"
            @save="updating(item,editable,header.value);editable=value"
            @cancel="editable = null">
          <span v-bind:class="isUpdated(item.id,header.value) ? 'updated' : ''">{{ value || ' ' }}</span>
          <template v-slot:input>
            <v-textarea
                class="mt-3"
                rows="4"
                no-resize
                filled
                v-model="editable"
                dense
                hide-details
                single-line
            />
          </template>
        </v-edit-dialog>
      </template>
    </v-data-table>
    <div class="text-center pt-2">
      <v-pagination
          class="my-4"
          v-model="page"
          :length="pageCount"
          :total-visible="7"
      />
    </div>
    <v-dialog v-model="windowSave" width="500" persistent>
      <v-card>
        <v-alert type="warning" :value="true"> Внимание!
        </v-alert>
        <v-card-text class="text-md-center">
          <v-textarea
              readonly
              value="Внести изменения?"
              rows="1"
              flat
              solo
              auto-grow
              dense>
          </v-textarea>
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="windowSave = false; save()">
            Да
          </v-btn>
          <v-btn @click="windowSave = false;">
            Нет
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="windowAlert" width="500" persistent>
      <v-card>
        <v-alert type="warning" :value="true"> Внимание!
        </v-alert>
        <v-card-text class="text-md-center">
          <v-textarea
              readonly
              value="Отменить изменения?"
              rows="1"
              flat
              solo
              auto-grow
              dense>
          </v-textarea>
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="windowAlert = false; rollback()">
            Да
          </v-btn>
          <v-btn @click="windowAlert = false;">
            Нет
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="windowInportFile" width="500" persistent>
      <v-card>
        <v-alert type="" :value="true"> Импорт
        </v-alert>
        <v-card-text class="text-md-center">
          <v-file-input ref="file" label="Выбранный файл" dense outlined counter truncate-length="60"
                        v-model="file" accept="*.xlsx,*.xls,*.xlsb"/>
          <v-text-field
              type="number"
              no-resize
              label="Год"
              v-model="selectedYear"
              rows="1"
              dense>
          </v-text-field>
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="windowInportFile = false; submitFile()">
            Да
          </v-btn>
          <v-btn @click="windowInportFile = false;">
            Нет
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
<script>

import axios from "axios";
import {bus} from "@/main";

export default {
  data() {
    return {
      windowSave: false,
      windowAlert: false,
      windowInportFile: false,
      numberOfCopies: 1,


      selectedFile:null,
      page: 1,
      pageCount: 20,

      editable: null,
      postFilterValue: null,
      sensorFilterValue: [],

      headerTable: [
        {
          text: 'Действия',
          align: 'center',
          sortable: false,
          value: 'actions',
          divider: true,
          width: 150
        },
        {
          text: 'Пост',
          align: 'start',
          sortable: false,
          value: 'postId',
          divider: true,
          width: 300,
        },
        {
          text: 'Тип сенсора',
          align: 'start',
          sortable: false,
          value: 'sensorType',
          divider: true,
          width: 300,
        },
        {
          text: 'Дата начала',
          align: 'start',
          sortable: false,
          value: 'dateStart',
          divider: true,
          width: 275,
        },
        {
          text: 'Дата конца',
          align: 'start',
          sortable: false,
          value: 'dateEnd',
          divider: true,
          width: 275,
        },
        {
          text: 'Примечание',
          align: 'start',
          sortable: false,
          value: 'comment',
          divider: true,
          width: 300
        }
      ],
      exceptions: [],
      copyExceptions: [],
      posts: [],
      sensorTypes: [],
      rowsForDelete: [],
      changes: {
        updatingRows: [],
        addedRows: [],
        deletedRows: []
      },
      listUUID: [],
      menu: {
        menuItem: null,
        showMenu: false,
        x: 0,
        y: 0,
      },
      file:[],
      badDate: [],
      badRepeat: [],
      badPeriod: [],
      selectedYear: new Date().getFullYear(),
    }
  },
  computed: {
    filteredDesserts() {
      let conditions = [];
      if (this.postFilterValue) {
        conditions.push(this.filterPost);
      }
      if (this.sensorFilterValue.length > 0) {
        conditions.push(this.filterSensor);
      }
      if (conditions.length > 0) {
        return this.exceptions.filter((exception) => {
          return conditions.every((condition) => {
            return condition(exception);
          })
        })
      }
      return this.exceptions;
    },
    loadingData() {
      return !(this.posts.length > 0 && this.sensorTypes.length > 0 && this.exceptions.length > 0);
    },
  },
  methods: {
    submitFile(){
      const formData = new FormData()
      formData.append('file', this.file)
      formData.append('year', this.selectedYear)
      this.isLoading = true
      axios
          .post('/api/uploadFile', formData, {headers: {'Content-Type': 'multipart/form-data'}})
          .then((response) => {
            const exceptionWithId = response.data
            for (const exceptionWithIdElement of exceptionWithId) {
              exceptionWithIdElement.id = this.genUUID()
            }
            this.changes.addedRows= exceptionWithId
            this.exceptions.unshift(...exceptionWithId)})
          .catch(() => {
            this.isLoading = false
            this.message = "Не удалось загрузить файл, возможно он был изменён"
          })
    },
    itemClasses(item) {
      return {
        added: this.changes.addedRows.length > 0 ? this.changes.addedRows.filter(j => j.id === item.id).length > 0 : false,
        errorDate: this.badDate.length > 0 ? this.badDate.includes(item.id, 0) : false,
        errorRepeat: this.badRepeat.length > 0 ? this.badRepeat.includes(item.id, 0) : false,
        errorPeriod: this.badPeriod.length > 0 ? this.badPeriod.includes(item.id, 0) : false,
      }
    },
    validateDate(item) {
      if (!(item.dateStart === '' || item.dateEnd === '')) {
        if (new Date(item.dateStart) >= new Date(item.dateEnd)) {
          return false
        }
      }
      return true
    },
    showMenu(e, item) {
      this.menu.menuItem = item
      e.preventDefault()
      this.menu.showMenu = false
      this.menu.x = e.clientX
      this.menu.y = e.clientY
      this.$nextTick(() => {
        this.menu.showMenu = true
      })
    },
    checkErrors() {
      this.badDate = []
      this.badRepeat = []
      this.badPeriod = []
      bus.$emit('message-reset', {})
      let flagError = false
      //Проверка добавленных строк на пустые поля и дубли
      addedRows: for (let i = 0; i < this.changes.addedRows.length; i++) {
        const row = this.changes.addedRows[i]
        //Проверка на дубли
        for (let j = 0; j < this.exceptions.length; j++) {
          const row2 = this.exceptions[j]
          if (row.id !== row2.id &&
              row.postId === row2.postId &&
              row.sensorType === row2.sensorType &&
              row.dateStart === row2.dateStart &&
              row.dateEnd === row2.dateEnd
          ) {
            //Удаляем дубли
            this.deleteRow(row.id)
            i--
            continue addedRows
          }
        }
        //Проверка на пустые поля
        for (const argumentsKey in row) {
          if ((row[argumentsKey] === null || row[argumentsKey] === '') && argumentsKey !== 'comment') {


            bus.$emit('message-add', {
              message: "Заполните пустые ячейки",
              color: "red",
              timeout: -1
            });
            flagError = true
            break;
          }
        }
      }
      //Проверка на пересечение периодов дат
      for (let i = 0; i < this.changes.addedRows.length; i++) {
        const changesRow = this.changes.addedRows[i]
        for (let j = 0; j < this.exceptions.length; j++) {
          const row2 = this.exceptions[j]
          if ((changesRow.id !== row2.id && changesRow.postId === row2.postId && changesRow.sensorType === row2.sensorType)
              && (new Date(changesRow.dateStart) < new Date(row2.dateEnd) && new Date(row2.dateStart) < new Date(changesRow.dateEnd) || changesRow.dateStart === row2.dateStart || changesRow.dateEnd === row2.dateEnd)) {
            if (!this.badPeriod.includes(changesRow.id, 0) && !this.badPeriod.includes(row2.id, 0)) {
              bus.$emit('message-add', {
                message: "Периоды дат пересекаются\n" +
                    "Пост: " + this.posts.filter(item => item.id === changesRow.postId)[0].name + "\n" +
                    "Тип сенсора: " + this.sensorTypes.filter(item => item.id === changesRow.sensorType)[0].name + "\n" +
                    changesRow.dateStart.replace('T', ' ') + " - " +
                    changesRow.dateEnd.replace('T', ' ') + "\n" +
                    row2.dateStart.replace('T', ' ') + " - " +
                    row2.dateEnd.replace('T', ' ') + "\n",
                color: "#28dcff",
                timeout: -1,
              });
              this.badPeriod.push(changesRow.id)
              flagError = true
            }
          }
        }
      }
      //Проверка обновленных строк на пустые поля, дубли и пересечение дат
      for (let i = 0; i < this.changes.updatingRows.length; i++) {
        const Updatedrow = this.changes.updatingRows[i].new
        //Проверка на пустые поля
        for (const argumentsKey in Updatedrow.new) {
          if ((Updatedrow[argumentsKey] === null || Updatedrow[argumentsKey] === '') && argumentsKey !== 'comment') {
            bus.$emit('message', {
              message: "Заполните пустые ячейки",
              color: "red",
              timeout: -1,
            });
            flagError = true
          }
        }
        for (let j = 0; j < this.exceptions.length; j++) {
          const row2 = this.exceptions[j]
          if (Updatedrow.id !== row2.id && Updatedrow.postId === row2.postId && Updatedrow.sensorType === row2.sensorType) {
            if (Updatedrow.dateStart === row2.dateStart && Updatedrow.dateEnd === row2.dateEnd) {
              this.badRepeat.push(Updatedrow.id)
              this.badRepeat.push(row2.id)
              bus.$emit('message-add', {
                message: "Есть повторяющиеся строки \n" +
                    "Пост: " + this.posts.filter(item => item.id === Updatedrow.postId)[0].name + "\n" +
                    "Тип сенсора: " + this.sensorTypes.filter(item => item.id === Updatedrow.sensorType)[0].name + "\n" +
                    "Дата начала: " + Updatedrow.dateStart.replace('T', ' ') + "\n" +
                    "Дата конца: " + Updatedrow.dateEnd.replace('T', ' '),
                color: "#ffab00",
                timeout: -1
              });
              flagError = true
            } else if (new Date(Updatedrow.dateStart) < new Date(row2.dateEnd) && new Date(row2.dateStart) < new Date(Updatedrow.dateEnd) || Updatedrow.dateStart === row2.dateStart || Updatedrow.dateStart === row2.dateStart) {
              if (!this.badPeriod.includes(Updatedrow.id, 0) && !this.badPeriod.includes(row2.id, 0)) {
                bus.$emit('message-add', {
                  message: "Периоды дат пересекаются\n" +
                      "Пост: " + this.posts.filter(item => item.id === Updatedrow.postId)[0].name + "\n" +
                      "Тип сенсора: " + this.sensorTypes.filter(item => item.id === Updatedrow.sensorType)[0].name + "\n" +
                      Updatedrow.dateStart.replace('T', ' ') + " - " +
                      Updatedrow.dateEnd.replace('T', ' ') + "\n" +
                      row2.dateStart.replace('T', ' ') + " - " +
                      row2.dateEnd.replace('T', ' ') + "\n",
                  color: "#28dcff",
                  timeout: -1,
                });
                this.badPeriod.push(Updatedrow.id)
                flagError = true
              }
            }
          }
        }
      }
      //Проверка коректности дат
      for (const exception of this.exceptions) {
        if (!(exception.dateStart === '' || exception.dateEnd === '')) {
          if (new Date(exception.dateStart) >= new Date(exception.dateEnd) ||
              (new Date(exception.dateStart) < new Date('1753-01-01T00:00')) ||
              (new Date(exception.dateEnd) < new Date('1753-01-01T00:00'))) {
            this.badDate.push(exception.id)
            bus.$emit('message-add', {
              message: "Некоректно введена дата\n" +
                  "Пост: " + this.posts.filter(item => item.id === exception.postId)[0].name + "\n" +
                  "Тип сенсора: " + this.sensorTypes.filter(item => item.id === exception.sensorType)[0].name + "\n" +
                  "Дата начала: " + exception.dateStart.replace('T', ' ') + "\n" +
                  "Дата конца: " + exception.dateEnd.replace('T', ' '),
              color: "red",
              timeout: -1
            });
            flagError = true
          }
        }
      }
      return flagError
    },
    save() {
      axios
          .post("/api/saveExceptions", {
            updatingRows: this.changes.updatingRows.map(item => item.new),
            addedRows: this.changes.addedRows.map(row => {
              return {
                postId: row.postId,
                sensorType: row.sensorType,
                dateStart: row.dateStart,
                dateEnd: row.dateEnd,
                comment: row.comment
              }
            }),
            deletedRows: this.changes.deletedRows
          })
          .then(() => {
            this.reset()
            bus.$emit('message', {
              message: "Данные успешно сохранены",
              color: "green"
            });
            setTimeout(() => this.getExceptions(), 10)
          })
          .catch(() => {
            bus.$emit('message', {
              message: "Ошибка сохранения данных",
              color: "red"
            });
          });
    },
    reset() {
      this.rowsForDelete = []
      this.badDate = []
      this.badRepeat = []
      this.badPeriod = []
      this.listUUID = []
      this.changes = {
        updatingRows: [],
        addedRows: [],
        deletedRows: []
      }
    },
    rollback() {
      bus.$emit("message-reset")
      this.reset()
      this.exceptions = structuredClone(this.copyExceptions)
    },
    filterPost(value) {
      return value.postId === this.postFilterValue||value.postId===null
    },
    filterSensor(value) {
      return this.sensorFilterValue.includes(value.sensorType)||value.sensorType===null
    },
    getRow(id) {
      for (let i = 0; i < this.exceptions.length; i++) {
        if (this.exceptions[i].id === id) {
          return {row: this.exceptions[i], index: i}
        }
      }
    },
    multiDelete() {
      for (const row of this.rowsForDelete) {
        this.deleteRow(row.id)
      }
      this.rowsForDelete = []
    },
    deleteRow(itemId) {
      const row = this.getRow(itemId).row
      this.exceptions.splice(this.getRow(itemId).index, 1)
      //Если в начале id находится #, то это добавленная стока
      if (itemId[0] === '#') {
        for (let i = 0; i < this.changes.addedRows.length; i++) {
          if (this.changes.addedRows[i].id === itemId) {
            this.changes.addedRows.splice(i, 1)
            this.listUUID.splice(this.listUUID.indexOf(itemId.slice(1)), 1)
            return
          }
        }
      }
      for (let i = 0; i < this.changes.updatingRows.length; i++) {
        let row = this.changes.updatingRows[i]
        if (row.new.id === itemId) {
          this.changes.deletedRows.push(row.old)
          this.changes.updatingRows.splice(i, 1)
          return;
        }
      }
      this.changes.deletedRows.push(row)
    },
    addRow() {
      let newRow = {id: this.genUUID(), postId: null, sensorType: null, dateStart: '', dateEnd: '', comment: ''}
      this.exceptions.unshift(newRow)
      this.changes.addedRows.push(newRow)
    },
    cloningRow(item, number) {
      for (let i = 0; i < number; i++) {
        let newRow = Object.assign({}, item)
        newRow.id = this.genUUID()
        this.exceptions.unshift(newRow)
        this.changes.addedRows.push(newRow)
      }
      this.numberOfCopies = 1
      this.menu.showMenu = false
    },
    isUpdated(rowId, col) {
      for (const row of this.changes.updatingRows) {
        if (rowId === row.old.id) {
          return row.old[col] !== row.new[col];
        }
      }
      return false
    },
    genUUID() {
      let UUID = Math.floor(Math.random() * 900000 + 1000)
      if (this.listUUID.includes(UUID)) {
        UUID = this.genUUID()
      } else {
        this.listUUID.push(UUID)
      }
      return '#' + UUID
    },
    updating(item, value, column) {
      //Если строка не поменлясь, то ничего не делаем
      if (item[`${column}`] === value) {
        return;
      }
      // Если в строчке id начинается на #, то она была добавлена
      if (item.id[0] === '#') {
        item[column] = value
        return;
      }
      // Если строка уже обновлялась, то только обновлем из списка обновленных
      for (let i = 0; i < this.changes.updatingRows.length; i++) {
        const updatingRow = this.changes.updatingRows[i]
        if (item.id === updatingRow.old.id) {
          //Если строчка обновление совпадает со старой версией,
          //то удаляем ее из списка обновленных
          if (this.isEqual(item, updatingRow.old)) {
            this.changes.updatingRows.splice(i, 1)
            return;
          }
          updatingRow.new[column] = value
          return;
        }
      }
      Object.pr
      // Если строка до этого не менялась, то добавляем ее версию до изменения old и новую версию new
      const oldItem = Object.assign({}, item);
      item[column] = value
      this.changes.updatingRows.push({old: oldItem, new: item})
    },
    isEqual(obj1, obj2) {
      for (const Key in obj1) {
        if (obj1[Key] !== obj2[Key]) {
          return false
        }
      }
      return true
    },
    getSensorTypes() {
      axios
          .post("/api/getSensorTypes", {})
          .then(response => {
            this.sensorTypes = response.data
          })
          .catch(e => {
            bus.$emit('message', {
              message: "Ошибка при получении сенсоров",
              color: "red"
            });
          });
    },
    getExceptions() {
      axios
          .post("/api/getExceptions", {})
          .then(response => {
            this.exceptions = response.data
            this.copyExceptions = structuredClone(response.data)
          })
          .catch(e => {
            bus.$emit('message', {
              message: "Ошибка при получении исключений",
              color: "red"
            });
          });

    },
    getPosts() {
      axios
          .post("/api/getPosts", {})
          .then(response => {
            this.posts = response.data
          })
          .catch(e => {
            bus.$emit('message', {
              message: "Ошибка при получении постов",
              color: "red"
            });
          });

    },
    /*    logout() {
          axios
              .post("/api/auth/logout", {})
              .then(response => {
                console.log(response.status)
              })
              .catch(() => {
                bus.$emit('message', {
                  message: "Ошибка во время выхода. Перезагрузите стринцу и попробуйте снова.",
                  color: "red"
                });
              });
        }*/
  },
  created() {
    bus.$emit('message-reset')
    this.getExceptions()
    this.getPosts()
    this.getSensorTypes()
  }
}
</script>
<style scoped>

.v-data-table >>> .v-text-field > .v-input__control > .v-input__slot:before, .v-text-field > .v-input__control > .v-input__slot:after {
  opacity: 0 !important;
}

.v-data-table >>> .theme--light.v-text-field:not(.v-input--has-state):hover > .v-input__control > .v-input__slot:before {
  opacity: 0 !important;
}

/*.v-data-table >>> .updated {
  background-color: #f4ffc4 !important;
}*/
.v-data-table >>> tr.added {
  background-color: #f4ffc4 !important;
}

.v-data-table >>> tr.errorDate {
  background-color: rgba(229, 56, 56, 0.55) !important;
}

.v-data-table >>> tr.errorRepeat {
  background-color: #ffab00 !important;
}

.v-data-table >>> tr.errorPeriod {
  background-color: #28dcff !important;
}


.v-data-table >>> tr:hover:not(.v-data-table__expanded__content):not(.v-data-table__empty-wrapper).added {
  filter: brightness(0.95) !important;
}

/*.v-data-table >>> tr:hover:not(.v-data-table__expanded__content):not(.v-data-table__empty-wrapper) .updated {
  filter: brightness(0.95) !important;
}*/

/*.v-data-table >>> tr:hover:not(.v-data-table__expanded__content):not(.v-data-table__empty-wrapper).added .updated {
  filter: inherit !important;
}*/

.date {
  margin: auto;
}

.date .v-icon {
  position: absolute;
  margin-left: -32px;
}

.date >>> .v-input__append-inner {
  padding: 0 !important;
}

.cell {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: nowrap;
  flex-direction: row;
}

.update-dot{
  margin-top: 3px;
}
.error-dot{
  margin-top: 3px;
}

.v-list {
  height: 350px; /* or any height you want */
  overflow-y: auto
}

/*th {
  vertical-align: middle;
  text-align: left;
  padding: 0 16px;
  user-select: none;
  font-size: 0.75rem;
  height: 32px;
  color: rgba(0, 0, 0, 0.6);
  border-bottom: thin solid rgba(0, 0, 0, 0.12);
  border-right: thin solid rgba(0, 0, 0, 0.12);
}*/
</style>