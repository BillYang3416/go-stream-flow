import { Component } from '@angular/core';
import { environment } from 'src/environments/environment.dev';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  lineLoginUrl = environment.apiUrl + '/auth/line-login';
}
