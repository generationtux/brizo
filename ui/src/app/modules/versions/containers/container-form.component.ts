import { Component, Input, Output, EventEmitter } from '@angular/core';
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
  @Output() delete: EventEmitter<any> = new EventEmitter();

  private selectedVolumeMount: string;
  private newPortProtocol: string = "TCP";
  private newPortNumber: number;
  private newVolumeMountName: string;
  private newVolumeMountPath: string;

  private onDelete(index: number) {
    this.delete.emit(null);
  }

  private addPort(e: any) {
    e.preventDefault();
    this.container.ports.push(new ContainerPort({
      protocol: this.newPortProtocol,
      port: this.newPortNumber,
    }));

    this.newPortNumber = null;
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

  private addVolumeMount(e: any) {
    e.preventDefault();

    this.container.volumeMounts.push(new VolumeMount({
      name: this.newVolumeMountName,
      path: this.newVolumeMountPath,
    }));

    this.newVolumeMountPath = "";
  }

  private removeVolumeMount(i: number) {
    this.container.volumeMounts.splice(i, 1);
  }
}
