import { Component, Input } from '@angular/core'
import { FormGroup } from '@angular/forms';
import { Volume } from '../version.model';

@Component({
    selector: 'volumeForm',
    templateUrl: 'volume-create.html',
})
export class VolumeCreateComponent {
    @Input() volume: Volume;
}
