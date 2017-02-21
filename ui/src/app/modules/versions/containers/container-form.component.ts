import { Component, Input } from '@angular/core';
import { Volume } from '../version.model';
import { Container, ContainerPort, VolumeMount } from './container.model';

@Component({
  selector: 'container-form',
  templateUrl: './container-form.html',
})

export class ContainerFormComponent {
  @Input() container: Container = new Container();
  @Input() containerIndex: number;
  @Input() availableVolumes: Volume[];

  private selectedVolumeMount: string;

  private addPort() {
    this.container.ports.push(new ContainerPort());
  }

  private removePort(i: number) {
    this.container.ports.splice(i, 1);
  }

  private addArg() {
    this.container.args.push("");
  }

  private removeArg(i: number) {
    this.container.args.splice(i, 1);
  }

  private addVolumeMount() {
    let mount = new VolumeMount({
      name: this.selectedVolumeMount,
    });
    this.container.volumeMounts.push(mount);
  }

  private removeVolumeMount(i: number) {
    this.container.volumeMounts.splice(i, 1);
  }
}
