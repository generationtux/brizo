import { Component, Input } from '@angular/core'
import { FormGroup } from '@angular/forms';

@Component({
    selector: 'volume',
    templateUrl: 'volume-create.html',
})
export class VolumeCreateComponent {
    @Input('group')
    public volumeForm: FormGroup;
}
