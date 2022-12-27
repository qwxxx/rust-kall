<template>
  <v-app class="app">
    <v-container>
      <div class="overflow position-fixed"
           v-if="isAlertVisible"
      >
        <v-alert class="alert"
                 dense
                 text
                 v-bind:type="isFetchSuccess ? 'success' : 'error'"
        >
          {{alertText}}
        </v-alert>
      </div>

      <v-form
          class="py-5"
          v-if="form === Form.LOGIN"
      >

        <v-text-field
            label="Пароль"
            type="password"
            v-model="password"
        ></v-text-field>

        <v-btn
            block
            color="primary"
            v-on:click="login"
        >
          Войти
        </v-btn>

      </v-form>

      <div
          class="pt-5 pb-10"
          v-else-if="form===Form.MAIN"
      >


        <div class="d-flex justify-space-between align-start">
          <v-col cols="6" class="pl-0">
            <v-select
                v-if="user_id===1"
                label="Анализ"
                v-model="analysisItemSelect"
                :items="analysisItems"
                @update:modelValue="analysisItemChange"
            ></v-select>
          </v-col>

          <v-col cols="2" class="pr-0 pl-0">
            <v-btn
                block
                color="primary"
                class="mb-5"
                v-on:click="logout()">
              <v-icon left >mdi-logout</v-icon>
              {{user_id==1?'Тренер':'Ученик'}}

            </v-btn>
          </v-col>
        </div>



        <div class="" >
          <div v-if="analysis===Analysis.TOURNAMENT">
            <div class="d-flex flex-column  justify-start">
              <div
                  class="d-flex justify-center align-center">
                <v-text-field
                    class="col-11 pt-0 mt-0"
                    label="ID турнира"
                    type="number"
                    v-model="tournamentId"
                    v-on:keyup.enter="findTournament()"
                ></v-text-field>

                <v-btn
                    color="primary"
                    class="col-2 ml-5 mb-9"
                    v-on:click="findTournament()">
                  <v-icon left>
                    mdi-magnify
                  </v-icon>
                  Найти

                </v-btn>
              </div>

            </div>

            <v-row class="flex-column-reverse">
              <v-col cols="12" class=" mb-5 mt-5"
                     v-for="[idx,tournament] in tournamentData">
                <v-table density="compact">
                  <template v-slot:default>
                    <thead >
                    <tr>

                      <th colspan="3"  style="font-size:15px;">
                        <v-row >
                          <v-col cols="2" class=""><b><br>{{"$"+tournament.stake}}</b></v-col>
                          <v-col cols="8" class="text-center"><br><b class="text-center">{{idx}}</b></v-col>
                          <v-col cols="2" class="d-flex justify-end">
                            <v-btn color="primary" class="mt-3 mr-2" v-on:click="getTournament(idx)">
                              <v-icon>mdi-sync</v-icon>
                            </v-btn>
                            <v-btn
                                color="error"
                                class="mt-3"
                                v-on:click="removeTournamentById(idx)"
                            >
                              <v-icon>
                                mdi-close
                              </v-icon>
                            </v-btn>
                          </v-col>
                        </v-row>

                      </th>
                    </tr>
                    <tr>
                      <th>Место</th>
                      <th class="text-center" >
                        Имя
                      </th>
                      <th class="text-center"  >
                        Очки
                      </th>
                    </tr>
                    </thead>

                    <tbody></tbody>
                    <tbody>
                    <tr
                        v-for="(player,placenum) in tournament.players"
                        v-bind:key="placenum"
                    >
                      <td>
                        {{JSON.stringify(placenum+1)}}
                      </td>
                      <td class="text-center">
                        <div class="d-flex flex-row justify-center"><v-icon
                            color="primary"
                            v-if='player.blocked'
                        >
                          mdi-lock
                        </v-icon>
                          <div >
                            {{player.name}}
                          </div>
                        </div>
                      </td>
                      <td class="text-center" >
                        <span v-if="player.blocked&&player.name==''">-</span>
                        <span v-else-if="player.unknown">{{player.score}}(нет в базе)</span>
                        <span v-else>{{player.score}}</span>

                      </td>
                    </tr>
                    <tr style="background-color:#21719d" >
                      <td><b>Всего</b></td>
                      <td></td>
                      <td class="text-center"><b>{{tournament.total_score}}</b></td>
                    </tr>
                    </tbody>
                  </template>
                </v-table>
              </v-col>
            </v-row>

            <br>
            <br>
          </div>
        </div>


        <div class="" v-if="analysis === Analysis.PLAYER">

          <div class="d-flex flex-column  justify-start">

            <div class="d-flex justify-center align-center">
              <v-text-field
                  class="col-11 pt-0 mt-0"
                  label="Никнейм игрока"
                  type="text"
                  v-model="playerName"
              ></v-text-field>

              <v-btn
                  color="primary"
                  class="col-2 ml-5 mb-9"
                  v-on:click="findPlayerInfo()">
                <v-icon left>
                  mdi-magnify
                </v-icon>
                Найти

              </v-btn>

            </div>

            <h3 class="mb-5">
              Период проведения турниров
            </h3>
            <v-row cols="12">


              <v-col class="col-6">
                <p class="mb-1">Начальная дата</p>
                <v-row>
                  <v-col cols="3" class="pr-1">
                    <v-text-field
                        id="startDayInput"
                        class="pt-0 mt-0"
                        v-model="startDay"
                        label="День"
                        type="number"
                        :rules="[isValidStartDay]"
                    ></v-text-field>
                  </v-col>
                  <v-col class="pr-1 pl-1">
                    <v-select
                        id="startMonthInput"
                        class="pt-0 mt-0"
                        v-model="startMonth"
                        :items="itemsMonth"
                        label="Месяц"
                        type="text"

                        max-length="2"
                    ></v-select>
                  </v-col>
                  <v-col cols="3" class="pl-1">
                    <v-text-field
                        id="startYearInput"
                        class="pt-0 mt-0"
                        v-model="startYear"
                        label="Год"
                        type="number"
                        :rules="[isValidStartYear]"
                    ></v-text-field>
                  </v-col>
                </v-row>
              </v-col>

              <v-col class="col-6">
                <p class="mb-1">Конечная дата</p>
                <v-row>
                  <v-col cols="3" class="pr-1">
                <v-text-field
                    id="endDayInput"
                    class="pt-0 mt-0"
                    v-model="endDay"
                    label="День"
                    type="number"
                    :rules="[isValidEndDay]"
                ></v-text-field>
                  </v-col>
                  <v-col class="pr-1 pl-1">
                <v-select
                    id="endMonthInput"
                    class="pt-0 mt-0"
                    v-model="endMonth"
                    :items="itemsMonth"
                    label="Месяц"
                    type="text"

                    max-length="2"
                ></v-select>
                  </v-col>
                  <v-col cols="3" class="pl-1">
                <v-text-field
                    id="endYearInput"
                    class="pt-0 mt-0"
                    v-model="endYear"
                    label="Год"
                    type="number"
                    :rules="[isValidEndYear]"
                ></v-text-field>
                  </v-col>
                </v-row>
              </v-col>


            </v-row>

          </div>

          <v-dialog
              id="1234"
              class="dialog__info-player"
              v-model="dialog"
              max-width="600"
          >
            <v-card>
              <v-card-title class="text-h5">
                Имя игрока: {{playerName}}
              </v-card-title>

              <v-card-text>
                Надено турниров  с {{startDate}} по {{endDate}} : <strong>{{tournamentsCount}}</strong>
                <br>
                Средний результат : <strong>{{averageTournamentsScore}}</strong>
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>

                <v-btn
                    color="green darken-1"
                    text
                    @click="dialog = false"
                >
                  Закрыть
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>

        </div>




          <div class=""
            v-if="user_id==1">

            <div class="d-flex justify-space-between">
                <div class="d-flex flex-column  justify-start col-5">
                <div class="text-center mb-5">Настройка БД</div>
                <div class="d-flex flex-column justify-space-between">

                    <v-btn class="mb-5 "
                        color="success"  download="db.xlsx" :href="'/api/config?password='+(password)" >

                    <v-icon left>mdi-download</v-icon>
                    Скачать базу
                    </v-btn>
                    <v-btn class="mb-5"
                        color="primary" v-on:click="chooseFiles()">
                    <v-icon left>mdi-file</v-icon>
                    {{dataFile===null?"Выбрать файл":"Выбрать другой файл"}}
                    </v-btn>
                    <v-btn class="mb-5"
                        color="primary" v-on:click="sendFile()"
                        v-if="dataFile!==null">
                    <v-icon left>mdi-upload</v-icon>
                    Отправить {{getUploadingFilesLabel()}}
                    </v-btn>
                    <input style="display:none" v-on:change="filesPicked" id="fileinput" type="file" accept=".xlsx"/>

                </div>

                <div class="d-flex flex-column align-start mt-5"
                v-if="analysis === Analysis.PLAYER">
                    <div class="bot-state mb-4">
                        <v-icon>
                            {{botState.icon}}
                        </v-icon>
                        {{botState.state}}
                        </div>
                        <v-btn color="primary"
                        @click="restartSystem">
                            <v-icon>mdi-reload</v-icon>
                            Перезагрузить систему
                        </v-btn>
                    </div>
                </div>

                <div class=" d-flex flex-column justify-start col-5">
                <div class="text-center mb-5">Список неизвестных имен</div>
                <div class=" d-flex flex-column justify-space-between ">

                    <v-btn
                        class="mb-5"
                        color="success"  download="unknownNames.xlsx" :href="'/api/unknownNames?password='+(password)">

                    <v-icon left>mdi-download</v-icon>
                    Скачать список
                    </v-btn>


                    <v-btn
                        color="error"  >

                    <v-icon left>mdi-cancel</v-icon>
                    Очистить



                    <v-dialog
                        v-model="dialogAccept"
                        activator="parent"
                    >
                        <v-card>
                        <v-card-title class="text-h5">
                            Вы уверены?
                        </v-card-title>
                        <v-card-text></v-card-text>
                        <v-card-actions>
                            <v-spacer></v-spacer>
                            <v-btn
                                color="green darken-1"
                                text
                                @click="dialogAccept = false"
                            >
                            Нет
                            </v-btn>
                            <v-btn
                                v-on:click="clearUnknownNames()"
                                color="green darken-1"
                                text
                                @click="dialogAccept = false"
                            >
                            Да
                            </v-btn>
                        </v-card-actions>
                        </v-card>
                    </v-dialog>
                    </v-btn>

                </div>
            </div>
          </div>


          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>

          <div class="d-flex" >
            <v-col cols="11"></v-col>

            <v-col cols="2" class="pr-0 pl-0 text-center d-flex flex-column" style="font-size:10px; width:10vw">
              <img style="width:10vw" src="./assets/logof.png" class="mx-auto d-block">
              <span >by ANDREIIANN</span>
            </v-col>
          </div>

        </div>
      </div>


      <div id="dialogbox" class="slit-in-vertical">
        <div>
          <div id="dialogboxbody"></div>
        </div>
      </div>

    </v-container>
  </v-app>
</template>

<script>
import { isProxy, toRaw } from 'vue'
import {ref} from '@vue/reactivity'

const Form = {
  LOGIN: 0,
  MAIN: 1
}
const UserID={
  UNKNOWN:0,
  ADMIN: 1,
  USER: 2
}
const Analysis = {
  TOURNAMENT: 0,
  PLAYER: 1
}

export default {
  data: () => ({
    alertText:"",
    form: Form.LOGIN,
    Form,
    analysis: Analysis.TOURNAMENT,
    Analysis,
    user_id:UserID.UNKNOWN,
    password: '',
    botState: {
      state: 'Состояние бота',
      icon: 'mdi-android'
    },
    dataFile:null,

    isAlertVisible: false,
    isAddingFilesVisible: false,
    isFetchSuccess: false,
    isFetchInProcess: false,


    dialog:false,
    dialogAccept:false,

    tournamentId:null,

    cachedTournamentIds:[28529378,28529380,28529406,28529356],
    notFoundTournamentsIds:[],
    errorTournamentsIds:[],
    tournamentData:new Map(),

    analysisItems: [
      'Анализ турнира',
      'Анализ игрока'
    ],
    analysisItemSelect: 'Анализ турнира',

    playerName: '',
    averageTournamentsScore: 0,
    tournamentsCount:0,
    startDate: '',
    endDate: '',

    itemsMonth: [...Array(12).keys()].map( key => new Date(0, key).toLocaleString('ru', { month: 'long' }) ),
    startDay: '2',
    startMonth: 'февраль',
    startYear: '2021',
    endDay: '',
    endMonth: '',
    endYear: '',

    deleteIcon: "mdi-delete",
    plusIcon: "mdi-plus-box-multiple",
  }),

  async mounted() {
    // console.log('mount')
    // console.log(import.meta.env.VITE_VUE_APP_HOST)
    // console.log(import.meta.env.VITE_VUE_APP_SERVER)
    const password = localStorage.getItem('password')

    if(!password) {
      return;
    }

    this.password = password

    let response=await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }login?password=${this.password}`, {
      method: "post",
    }).catch((error) => ("Something went wrong!", error));
    if(response.status>=400){
      return
    }
    var r=await response.json()
    this.user_id=r.user_id

    window.setInterval(() => {
      this.updateTournament()
    }, 5000)
    await this.toMain()

    let cdate=new Date();
    this.endDay=cdate.getDate();
    this.endMonth=this.itemsMonth[cdate.getMonth()]
    this.endYear=cdate.getFullYear()

    this.getBotState()
    window.setInterval( () => {
      this.getBotState()
    }, 10000)
  },

  methods: {
    async restartSystem() {
      const response = await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }restart`) 
      if (response.status === 200) {
          return
      } else {
          this.customAlert('Ошибка')
      }
    },
    async getBotState() {
      const response = await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }state`)    
      if (response.status === 200) {
        const data = response.json()
        if (data.isPlayerCalculateRunning === true) {
          this.botState = {
            state: 'Идёт подсчёт',
            icon: 'mdi-progress-clock'
          }
        } else {
            this.botState = {
            state: 'Бот готов',
            icon: 'mdi-android'
          }
        }
      } else {
        this.customAlert('Ошибка')
      }
    },
    validateDay(day) {
      let n = `${day}`
      if(n.length >= 2) {
        n = n.slice(0,2)
        if (n > 31) n = `31`
        if (n < 1) n = `1`
      }
      return n;
    },
    isValidStartDay(a) {
      a=this.validateDay(a)
      this.startDay= a
      return true
    },
    isValidEndDay(a) {
      a=this.validateDay(a)
      this.endDay= a
      return true
    },
    validateYear(year) {
      let n = `${year}`
      if(n.length >=4) {
        n = n.slice(0,4)
        if( n > (new Date()).getFullYear() ){
          n = (new Date()).getFullYear();
        }else if(n < '2000'){
          n = '2000'
        }
      }
      return n;
    },
    isValidStartYear(a) {
      a = this.validateYear(a)
      this.startYear = a
      return true
    },
    isValidEndYear(a) {
      a = this.validateYear(a)
      this.endYear = a
      return true
    },

    removeTournamentById(idx){
      this.tournamentData.delete(idx)
    },
    getUploadingFilesLabel(){
      if(this.dataFile){

        return this.dataFile.name
      }else{
        return ''//'Файлы не выбраны'
      }
    },
    chooseFiles(){
      document.getElementById("fileinput").click()
    },
    async sendFile(){
      const formData = new FormData();


      formData.append("files",this.dataFile);
      for (var pair of formData.entries()) {
        console.log(pair[0]+ ', ' + pair[1]);
      }
      let response=await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }api/config?password=${this.password}`, {
        method: "post",
        body: formData,
      }).catch((error) => ("Something went wrong!", error));
      console.log(response)
      function onlyUnique(value, index, self) {
        return self.indexOf(value) === index;
      }

      if(response.status===400){
        let repNames=await response.json()
        repNames=repNames.filter(onlyUnique)
        let str="Найдены повторяющиеся имена: "
        for (let name of repNames){
          str+=name
          str+=","
        }
        str=str.substring(0,str.length-1)
        this.customAlert(str)
      }
      else if(response.status>400){
        this.customAlert("Неизвестная ошибка!")
        return;
      }else{
        this.customAlert('Файл успешно загружен! Повторяющихся имен не обнаружено.')
      }

      this.dataFile=null
      document.getElementById('fileinput').value=null;
    },
    filesPicked(){
      console.log('picked files')
      let inputFile=document.getElementById('fileinput');
      this.dataFile=inputFile.files[0]
    },
    async loadFileClicked(){

    },
    async updateTournament(){
      console.log(this.tournamentData)
      for(let [tid,tval] of this.tournamentData){
        if (tval!=null&&tval.players[5].name==='6 место'){
          await this.getTournament(tid)
        }
      }

    },
    async getTournament(tid){
      if (tid===null||tid===""){
        return;
      }
      let response = await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }api/tournament?network=WPN&password=${this.password}&tournament_id=${tid}`)

      if (response.status>=400) {
        if (response.status == 404) {
          this.customAlert(`Турнир с ID ${tid} не найден`)
        }else{
          this.customAlert(`Ошибка при поиске турнира с ID ${tid}`)
        }
      }else{
        this.tournamentData.set(tid,await response.json())
      }

    },
    async getTournaments(){
      if(this.cachedTournamentIds==null){
        return
      }
      console.log(this.cachedTournamentIds)
      this.tournamentData=new Map()
      for(let tid of this.cachedTournamentIds){
        await this.getTournament(tid)
      }
      if(this.notFoundTournamentsIds.length>0){
        let str="Турниры со следующими номерами не найдены: "
        for(let id of this.notFoundTournamentsIds){
          str+=id
          str+=","
        }
        str=str.substring(0,str.length-1)
        this.customAlert(str)


      }
      if(this.errorTournamentsIds.length>0){
        let str="Ошибка поиска для турниров со следующими номерами: "
        for(let id of this.errorTournamentsIds){
          str+=id
          str+=","
        }
        str=str.substring(0,str.length-1)
        this.customAlert(str)


      }

      this.errorTournamentsIds=[]
      this.notFoundTournamentsIds=[]


    },
    async findTournament(){
      console.log('find')
      console.log(this.tournamentId)
      this.cachedTournamentIds.push(this.tournamentId)
      await this.getTournament(this.tournamentId)
    },
    async clearUnknownNames(){
      let response=await fetch(
          `${ import.meta.env.VITE_VUE_APP_SERVER }api/unknownNames/clear?password=${this.password}`,
          {
            method:"post",
          })

      if(response.status>=400){
        this.customAlert('Произошла ошибка при очистке!')
      }else{
        this.customAlert('Очищено успешно!')
      }


    },
    logout(){
      localStorage.clear()
      this.form=Form.LOGIN
      window.location.reload()
    },
    async login(){
      console.log('login')
      let response=await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }login?password=${this.password}`, {
        method: "post",
      }).catch((error) => ("Something went wrong!", error));
      console.log(response)
      if(response.status>=400){
        this.isAlertVisible=true;
        this.isFetchSuccess=false;
        this.alertText='Неверный пароль'
        setTimeout(() => {
          this.isAlertVisible=false;
        }, 2000);
        return;
      }
      var r=await response.json()
      console.log(r)
      this.user_id=r.user_id
      console.log(this.user_id)

      localStorage.setItem('password', this.password)
      await this.toMain();
    },
    async switchFilterActive(){
      this.filter.is_active=!this.filter.is_active;
    },
    async toMain(){
      this.form=Form.MAIN
    },

    analysisItemChange() {
      if (this.analysisItemSelect === 'Анализ игрока') {
        this.analysis = this.Analysis.PLAYER
      }
      if (this.analysisItemSelect === 'Анализ турнира') {
        this.analysis = this.Analysis.TOURNAMENT
      }
    },

    async findPlayerInfo() {
      if (!this.validatePlayerInputs()) return

      const response = await fetch(`${ import.meta.env.VITE_VUE_APP_SERVER }api/player?network=WPN&playerName=${this.playerName}&startDate=${this.startDate}&endDate=${this.endDate}&password=${this.password}`)

      console.log("ALERT: ",response)
      let data = await response.json()



      if (response.status === 500) {
        this.customAlert('Непредвиденная ошибка поиска')
        return
      }
      if (response.status>=400) {
        if (response.status == 404) {
          this.customAlert(`Игрок с ником ${playerName} не найден`)
          return
        }else{
          this.customAlert(`Ошибка при поиске игрока с ником ${playerName}`)
          return
        }
      }
      console.log("ALERT: "+data.error)
      if (data.error === 0){
        this.playerName = data.playerName
        this.averageTournamentsScore = data.average_tournaments_score
        this.averageTournamentsScore = this.averageTournamentsScore.toFixed(2)
        this.tournamentsCount = data.tournaments_count
        this.startDate = new Date(data.start_date*1000).toLocaleDateString('ru-RU')
        this.endDate = new Date(data.end_date*1000).toLocaleDateString('ru-RU')
        this.dialog = true
      }

      if (data.error === 1) {
        this.customAlert('Игрок не найден')
      }
      if (data.error === 2) {
        this.customAlert('Турниров не найдено')
      }

    },

    isValidDate(day, month, year) {
      let dateString = `${day}/${month}/${year}`
      // First check for the pattern
      if(!/^\d{1,2}\/\d{1,2}\/\d{4}$/.test(dateString))
      {
        console.log(false)
        return false;
      } else {
        console.log(true)
      }

      // Parse the date parts to integers
      var parts = dateString.split("/");
      var day = parseInt(parts[0], 10);
      var month = parseInt(parts[1], 10);
      var year = parseInt(parts[2], 10);

      console.log(parts, day, month, year)

      // Check the ranges of month and year
      if(year < 1970) {
        console.log(false)
        return false
      }
      if(year > 3000) {
        console.log(false)
        return false
      }
      if(month == 0 || month > 12) {
        console.log(month)
        console.log(false)
        return false
      }

      var monthLength = [ 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31 ];

      // Adjust for leap years
      if(year % 400 == 0 || (year % 100 != 0 && year % 4 == 0)) monthLength[1] = 29;
      // Check the range of the day
      if (day > 0 && day <= monthLength[month - 1]) console.log(true);
      return day > 0 && day <= monthLength[month - 1];
    },

    validatePlayerInputs() {
      if (this.playerName === '') {
        this.customAlert('Введите никнейм игрока')
        return false
      }
      const currentStartMonth = this.itemsMonth.indexOf(`${this.startMonth}`) + 1
      let startDayTwoChars=`${this.startDay}`
      if(startDayTwoChars.length===1){
        startDayTwoChars='0'+startDayTwoChars
      }
      let startMonthTwoChars=`${currentStartMonth}`
      if (startMonthTwoChars.length===1){
        startMonthTwoChars='0'+startMonthTwoChars
      }
      console.log(`${this.startYear}-${startMonthTwoChars}-${startDayTwoChars}`)
      if (this.isValidDate(this.startDay,  currentStartMonth , this.startYear)) {
        const newFormatDate = `${this.startYear}-${startMonthTwoChars}-${startDayTwoChars}`
        const newdate = `${Math.floor(new Date(newFormatDate).getTime() / 1000)}`
        this.startDate = newdate
        console.log(this.startDate)
      } else {
        this.customAlert('Некорректная начальная дата')
        return false
      }
      let dayTwoChars=`${this.endDay}`;
      if (dayTwoChars.length===1){
        dayTwoChars='0'+dayTwoChars
      }
      let monthTwoChars=`${this.itemsMonth.indexOf(`${this.endMonth}`) + 1}`

      if(monthTwoChars.length===1){
        monthTwoChars=`0${monthTwoChars}`
      }
      console.log(`${this.endYear}-${monthTwoChars}-${dayTwoChars}`)
      if (this.isValidDate(this.endDay,  this.itemsMonth.indexOf(`${this.endMonth}`) + 1 ,this.endYear)) {
        const newFormatDate = `${this.endYear}-${monthTwoChars}-${dayTwoChars}`
        let newdate = Math.floor(new Date(newFormatDate).getTime() / 1000)
        newdate+=(60*60*24)
        this.endDate = `${newdate}`
      } else {
        this.customAlert('Некорректная конечная дата')
        return false
      }

      return true
    },

    customAlert(message) {
      const dialogbox = document.getElementById('dialogbox')

      dialogbox.style.top = "50px"
      dialogbox.style.right = "15px"

      dialogbox.style.display = "block"

      document.getElementById('dialogboxbody').textContent = message

      setTimeout(() => {
        dialogbox.style.animationName = 'slit-in-fade'
      }, 3400)
      setTimeout(() => {
        dialogbox.style.animationName = 'slit-in-vertical'
        dialogbox.style.display = "none"
      }, 3600)
    }
  },


}
</script>

<style>
.bot-state{
  border-radius: 5px;
  background-color: #4caf50;
  padding: 6px 12px;
}
td{
  width:33%;
}
th{
  width:33%;
}
.position-fixed {
  position: fixed;
}
.mnx-2 {
  margin-left: -8px;
  margin-right: -8px;
}

.v-text-field > .v-input__control > .v-input__slot > .v-text-field__slot {
  background-color: #333;
}

.app {
  background-color: #444 !important;
}

.v-input input {
  padding-left: 8px !important;
}

.v-input .v-label {
  padding-left: 8px;
}

.v-textarea textarea {
  min-height: 92px !important;
  padding-left: 8px !important;
}

.overflow {
  background-color: white;
  right: 15px;
  width: 250px;
  max-width: 100%;
  z-index: 2;
  border-radius: 4px;
}

.alert {
  margin: 0 auto !important;
  width: 100%;
  max-width: 100%;
  background: rgb(var(--v-theme-error)) !important;
  color: rgb(var(--v-theme-on-error)) !important;
  --v-theme-overlay-multiplier: var(--v-theme-error-overlay-multiplier);
  position: relative;
  padding: 16px;
  overflow: hidden;
  --v-border-color: currentColor;
  border-radius: 4px;
}
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

/* Firefox */
input[type=number] {
  -moz-appearance: textfield;
}

.dialog__info-player {
  width: 600px;
}

#dialogbox{
  display: none;
  position: absolute;
  z-index: 10;
  top:0;
  right: 0;
  background: rgb(37, 37, 37);
  border-radius: 8px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.575);
  transition: 0.3s;
  min-width: 250px;
  max-width: 450px;
  border-left: 4px solid rgb(var(--v-theme-error));
  padding: 16px;
  word-wrap: break-word;
}

#dialogbox:hover {
  box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.911);
}

.slit-in-vertical {
  -webkit-animation: slit-in-vertical 0.45s ease-out both;
  animation: slit-in-vertical 0.45s ease-out both;
}

@keyframes slit-in-vertical {
  0% {
    -webkit-transform: translateX(800px);
    transform: translateX(800px);
    opacity: 0;
  }
  54% {
    -webkit-transform: translateX(160px);
    transform: translateX(160px);
    opacity: 1;
  }
  100% {
    -webkit-transform: translateX(0);
    transform: translateX(0);
  }
}
@keyframes slit-in-fade {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}

</style>