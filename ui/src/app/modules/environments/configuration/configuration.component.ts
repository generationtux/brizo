import 'rxjs/add/operator/switchMap'
import { Observable } from 'rxjs/Rx'
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit, AfterViewInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms'

import { EnvironmentService } from '../environment.service'
import { ConfigurationService } from './configuration.service'

@Component({
  selector: 'configuration',
  templateUrl: './configuration.html',
  providers: [ConfigurationService]
})

export class ConfigurationComponent {
  private environment: any = {};
  private application: any = {};
  private configuration: any = [];
  private createConfigForm: FormGroup;
  private lastIndex: number = 0;
  
  constructor(private configService: ConfigurationService,
              private environmentService: EnvironmentService,
              private route: ActivatedRoute) {}
  
  ngAfterViewChecked() {
    this.lastIndex = this.configuration.length;
  }
  
  ngOnInit() {
    this.createConfigForm = new FormGroup({
        name: new FormControl('', [<any>Validators.required]),
        value: new FormControl('', [<any>Validators.required])
    });
    this.route.params
      .switchMap(
        (params: Params) => this.environmentService.getEnvironment(params['uuid']),
      ).subscribe(
        data => {
          this.environment = data;
          this.application = this.environment.application;
        },
        err => console.error('There was an error: ' + err)
      );
    this.route.params
      .switchMap(
        (params: Params) => this.configService.getConfiguration(params['uuid']),
      ).subscribe(
        data => {
          this.configuration = data;
        },
        err => console.error('There was an error: ' + err)
      );
  }
  
  ngAfterViewInit() {
    var i = $('tr[class^="addr"]').length;
    $(".add-btn").click(function() {
      $('.addr' + i).html("<td><input name='key' type='text' placeholder='KEY' class='form-control' /></td><td><input name='value' type='text' placeholder='VALUE' class='form-control' /></td><td><button class='btn btn-danger' type='button'><i class='fa fa-times'></i></button></td>");
      $('.config-table').append("<tr class='addr" + (i + 1) + "'></tr>");
      i++;
    });
    $(".del-btn").click(function() {
      if(i > 1) {
        $(".addr" + (i-1)).html('');
        i--;
		  }
    });
  }
  
  createNewConfig() {
    var data = $('.config-form').serializeArray();
    var config = {};
    
    for(var i = 0; i < data.length; i++) {
      if(data[i].name == 'key') {
        var key = data[i].value;
      }
      
      if(data[i].name == 'value') {
        var val = data[i].value;
      }
      
      config[key] = val;
    }
    
    for(var key in config) {
      var name  = key;
      var value = config[key];
      this.configService.createConfiguration(name, value, this.environment.uuid).subscribe(
        response => console.log('configuration created'),
        err => console.error(err),
      )
    }
  }
}