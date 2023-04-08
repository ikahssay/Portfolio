package com.application.cs260final;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.ImageView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

public class Notifications extends AppCompatActivity {
    ImageView forward_button_top;
    ImageView forward_button_bottom;
    ImageView home_button;
    ImageView search_button;
    ImageView predictions_button;
    ImageView notification_button;
    ImageView profile_button;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_notifications);


        //1. CREATE AN ACTION LINKED TO THE FIRST FORWARD BUTTON
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        forward_button_top = (ImageView) findViewById(R.id.imageForwardButton);
        forward_button_top.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(Notifications.this, "See Detailed Post", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(Notifications.this, ShowDetailedPostScreen.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //2. CREATE AN ACTION LINKED TO THE SECOND FORWARD BUTTON
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        forward_button_bottom = (ImageView) findViewById(R.id.imageForwardButton1);
        forward_button_bottom.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(Notifications.this, "See Detailed Post", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(Notifications.this, ShowDetailedPostScreen.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });



        // ************************************* MENU BUTTON FUNCTIONALITIES BELOW **************************************** //

        //3. CREATE AN ACTION LINKED TO THE HOME BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        home_button = (ImageView) findViewById(R.id.imageHome);
        home_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(Notifications.this, "Going to Home Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                //Notification is about user downvoting a post -> sent to screen where post was shown as downvoted.
                Intent intent = new Intent(Notifications.this, HomeScreenActivity.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //4. CREATE AN ACTION LINKED TO THE SEARCH BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        search_button = (ImageView) findViewById(R.id.imageSearch);
        search_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(Notifications.this, "Going to Search Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(Notifications.this, Search.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //5. CREATE AN ACTION LINKED TO THE PREDITCION BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        predictions_button = (ImageView) findViewById(R.id.imageVote);
        predictions_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(Notifications.this, "Going to Predictions Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(Notifications.this, Predictions.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });


        //6. CREATE AN ACTION LINKED TO THE PROFILE BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        profile_button = (ImageView) findViewById(R.id.imageProfileButton);
        profile_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(Notifications.this, "Going to Profile Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(Notifications.this, Profile.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });


    }
}