import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { Router, ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'

import { Environment } from '../../environments/environment.model';
import { EnvironmentService } from '../../environments/environment.service';
import { Version, Volume, VolumeMount, ContainerPort } from '../version.model';
import { VersionService } from '../version.service';

@Component({
  selector:    'version-create',
  templateUrl: './version-create.html',
})

export class VersionCreateComponent implements OnInit {
  public version: Version;
  private application: any = {}
  private environment: any = {}

  constructor(
    private versionService: VersionService,
    private environmentService: EnvironmentService,
    private activeRoute: ActivatedRoute,
    private router: Router,
  ) {
    this.version = new Version({
      volumes: [],
      args: [],
      ports: [],
      volumeMounts: [],
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
    console.log(this.version);
    // let version = new Version({});
    //
    // this.versionService.createVersion(version).subscribe(
    //   () => (this.onCreateVersion()),
    //   err => console.error('There was an error: ' + err)
    // );
  }

  private onCreateVersion() {
    this.router.navigate(['/environments', this.environment.uuid]);
  }

  private addVolume() {
    this.version.volumes.push(new Volume());
  }

  private removeVolume(i: number) {
    this.version.volumes.splice(i, 1);
  }

  private addPort() {
    this.version.ports.push(new ContainerPort());
  }

  private removePort(i: number) {
    this.version.ports.splice(i, 1);
  }

  private addVolumeMount(name: string) {
    this.version.volumeMounts.push(new VolumeMount({name}));
  }
}
