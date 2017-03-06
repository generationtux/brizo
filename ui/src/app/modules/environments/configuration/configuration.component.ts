import 'rxjs/add/operator/switchMap'
import { Observable } from 'rxjs/Rx'
import { Router, ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, Input, OnInit, AfterViewInit } from '@angular/core'
import { FormArray, FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms'

import { Configuration } from './configuration.model';
import { Environment } from '../environment.model';
import { EnvironmentService } from '../environment.service'
import { ConfigurationService } from './configuration.service'

@Component({
  selector: 'configuration',
  templateUrl: './configuration.html',
  styleUrls: ['./configuration.style.css'],
  providers: [ConfigurationService]
})

export class ConfigurationComponent {
  private createConfigForm: FormGroup;
  private environment: any = {};
  private application: any = {};

  constructor(private configService: ConfigurationService,
              private formBuilder: FormBuilder,
              private environmentService: EnvironmentService,
              private route: ActivatedRoute,
              private router: Router) {}

  ngOnInit() {
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
    this.route.params.switchMap((params: Params) => {
      return this.configService.getConfigurations(params['uuid']);
    }).subscribe(
      data => this.initConfigurations(data),
      err => console.error('There was an error: ' + err)
    );
    this.createConfigForm = this.formBuilder.group({
        configurations: this.formBuilder.array([]),
    });
  }

  private initConfigurations(configurations: Configuration[]) {
    const controls = <FormArray>this.createConfigForm.controls['configurations'];
    configurations.map(configuration => {
      controls.push(this.buildConfiguration(
        configuration.uuid,
        configuration.name,
        configuration.value,
      ));
    })
  }

  private buildConfiguration(uuid: string = '', name: string = '', value: string = ''): FormGroup {
    return this.formBuilder.group({
      uuid:  [uuid],
      name:  [name,  Validators.required],
      value: [value, Validators.required],
    });
  }

  private addConfiguration() {
    const controls = <FormArray>this.createConfigForm.controls['configurations'];
    controls.push(this.buildConfiguration());
  }

  private removeConfiguration(i: number) {
    const controls = <FormArray>this.createConfigForm.controls['configurations'];
    const configuration = controls.value[i];
    if (configuration.uuid !== '' || configuration.uuid !== null) {
      this.configService.deleteConfiguration(
        this.environment.uuid,
        configuration.uuid
      ).subscribe(
        data => {},
        err => console.error(err),
      )
    }
    controls.removeAt(i);
  }

  private saveConfigurations() {
    const configurations = this.createConfigForm.controls['configurations'].value;
    configurations.map((configuration: Configuration, index: number) => {
      if (configuration.uuid === '' || configuration.uuid === null) {
        this.configService.createConfiguration(
          configuration.name,
          configuration.value,
          this.environment.uuid,
        ).subscribe(
          data => {},
          err => console.error(err),
        );
      }
    });
    this.router.navigate(['/environments', this.environment.uuid]);
  }
}
