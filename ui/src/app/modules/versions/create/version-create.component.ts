import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { Router, ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'

import { Environment } from '../../environments/environment.model';
import { EnvironmentService } from '../../environments/environment.service';
import { Version, Volume } from '../version.model';
import { Container } from '../containers/container.model';
import { VersionService } from '../version.service';

@Component({
  selector:    'version-create',
  templateUrl: './version-create.html',
})

export class VersionCreateComponent implements OnInit {
  public version: Version;
  private application: any = {}
  private environment: any = {}
  public error: string;

  private newVolumeType: string = "temp";
  private newVolumeName: string;

  constructor(
    private versionService: VersionService,
    private environmentService: EnvironmentService,
    private activeRoute: ActivatedRoute,
    private router: Router,
  ) {
    this.version = new Version({
      replicas: 1,
      containers: [new Container({name: 'container-1'})],
      volumes: [],
    });
  }

  ngOnInit() {
    this.activeRoute.params.switchMap((params: Params) => {
      return this.environmentService.getEnvironment(params['environment-uuid']);
    }).subscribe(
      data => this.onGetEnvironment(data),
      err => console.error('There was an error: ' + err)
    );
  }

  private onGetEnvironment(data: any) {
    this.environment = data;
    this.application = data.application || {};
  }

  private createVersion(e: any) {
    e.preventDefault();
    this.version.environment_uuid = this.environment.uuid;

    this.versionService.createVersion(this.version).subscribe(
      () => this.onCreateVersion(),
      err => this.displayError(err)
    );
  }

  private onCreateVersion() {
    this.router.navigate(['/environments', this.environment.uuid]);
  }

  private displayError(e: string) {
    this.error = 'There was an error: ' + e
  }

  private addContainer() {
    this.version.containers.push(new Container({
      name: 'container-' + (this.version.containers.length + 1)
    }));
  }

  private removeContainer(i: number) {
    this.version.containers.splice(i, 1);
  }

  private addVolume(e: any) {
    e.preventDefault();
    this.version.volumes.push(new Volume({
      name: this.newVolumeName,
      type: this.newVolumeType,
    }));

    this.newVolumeName = "";
  }

  private removeVolume(i: number) {
    this.version.volumes.splice(i, 1);
  }
}
