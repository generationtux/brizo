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
  styleUrls: ['./configuration.style.css'],
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
    $(document).ready(function() {
      /* @TODO remove inputs and remove from db */
      $('.del-btn').click(function() {
        var id = $(this).parent().parent().data('config-id');
        
        /* @TODO finish out deleteConfiguration part of service
        this.configService.deleteConfiguration(parseInt(id), this.environment.uuid).subscribe(
          response => console.log('configuration deleted'),
          err => console.error(err),
        );
        */
      });
      
      $('.add-btn').on('click', function() {
        var i = $('ul.config-entry').length;
        $('.configs').append(''+
          '<ul class="list-inline config-entry clearfix addr' + (i + 1) + '">'+
          '<li class="config-name"><input class="form-control" name="key" formControlName="name" placeholder="KEY" type="text" /></li>'+
          '<li class="config-value"><input class="form-control" name="value" formControlName="value" placeholder="VALUE" type="text" /></li>'+
          '<li class="config-del"><button class="btn btn-block btn-danger del-btn" type="button" value="addr0"><i class="fa fa-times"></i></button></li>'+
          '</ul>');
      });
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